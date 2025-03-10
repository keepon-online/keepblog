# springboot mybatis aop 读写分离

##  注解

- 定义两个注解

  ```java
  
  import java.lang.annotation.ElementType;
  import java.lang.annotation.Retention;
  import java.lang.annotation.RetentionPolicy;
  import java.lang.annotation.Target;
  
  /**
   * 主库可读写
   */
  @Target(ElementType.METHOD)
  @Retention(RetentionPolicy.RUNTIME)
  public @interface Master {
  }
  ```

  ```java
  import java.lang.annotation.ElementType;
  import java.lang.annotation.Retention;
  import java.lang.annotation.RetentionPolicy;
  import java.lang.annotation.Target;
  
  /**
   * 从库可读
   */
  @Target(ElementType.METHOD)
  @Retention(RetentionPolicy.RUNTIME)
  public @interface Slave {
  }
  ```

  ##  数据库配置

  ```yaml
  
  spring:
    datasource:
      master:
        jdbc-url: jdbc:mysql://192.168.1.22:3307/test
        username: root
        password: 123456
        driver-class-name: com.mysql.cj.jdbc.Driver
      slave1:
        jdbc-url: jdbc:mysql://192.168.1.22:3307/test
        username: root   # 只读账户
        password: 123456
        driver-class-name: com.mysql.cj.jdbc.Driver
      slave2:
        jdbc-url: jdbc:mysql://192.168.1.22:3307/test
        username: root   # 只读账户
        password: 123456
        driver-class-name: com.mysql.cj.jdbc.Driver
  		
  ```

  - 定义 数据源

    ```java
    
    public enum  DBTypeEnum {
        MASTER, SLAVE1, SLAVE2;
    }
    
    ```

    ##  mybatis 配置

    

    数据源配置

    ```java
    import org.springframework.beans.factory.annotation.Qualifier;
    import org.springframework.boot.context.properties.ConfigurationProperties;
    import org.springframework.boot.jdbc.DataSourceBuilder;
    import org.springframework.context.annotation.Bean;
    import org.springframework.context.annotation.Configuration;
    
    import javax.sql.DataSource;
    import java.util.HashMap;
    import java.util.Map;
    
    @Configuration
    public class DataSourceConfig {
        @Bean
        @ConfigurationProperties("spring.datasource.master")
        public DataSource masterDataSource() {
            return DataSourceBuilder.create().build();
        }
    
        @Bean
        @ConfigurationProperties("spring.datasource.slave1")
        public DataSource slave1DataSource() {
            return DataSourceBuilder.create().build();
        }
    
        @Bean
        @ConfigurationProperties("spring.datasource.slave2")
        public DataSource slave2DataSource() {
            return DataSourceBuilder.create().build();
        }
    
        @Bean
        public DataSource routingDataSource(@Qualifier("masterDataSource") DataSource masterDataSource,
                                              @Qualifier("slave1DataSource") DataSource slave1DataSource,
                                              @Qualifier("slave2DataSource") DataSource slave2DataSource) {
            Map<Object, Object> targetDataSources = new HashMap<>();
            targetDataSources.put(DBTypeEnum.MASTER, masterDataSource);
            targetDataSources.put(DBTypeEnum.SLAVE1, slave1DataSource);
            targetDataSources.put(DBTypeEnum.SLAVE2, slave2DataSource);
            RoutingDataSource routingDataSource = new RoutingDataSource();
            routingDataSource.setDefaultTargetDataSource(masterDataSource);
            routingDataSource.setTargetDataSources(targetDataSources);
            return routingDataSource;
        }
    
    }
    
    ```

    mybatis 配置

    ```java
    
    
    import org.apache.ibatis.session.SqlSessionFactory;
    import org.mybatis.spring.SqlSessionFactoryBean;
    import org.springframework.context.annotation.Bean;
    import org.springframework.context.annotation.Configuration;
    import org.springframework.core.io.support.PathMatchingResourcePatternResolver;
    import org.springframework.jdbc.datasource.DataSourceTransactionManager;
    import org.springframework.transaction.PlatformTransactionManager;
    import org.springframework.transaction.annotation.EnableTransactionManagement;
    
    import javax.annotation.Resource;
    import javax.sql.DataSource;
    
    @Configuration
    @EnableTransactionManagement
    public class MyBatisConfig {
    
        @Resource(name = "routingDataSource")
        private DataSource routingDataSource;
    
        @Bean
        public SqlSessionFactory sqlSessionFactory() throws Exception {
            SqlSessionFactoryBean sqlSessionFactoryBean = new SqlSessionFactoryBean();
            sqlSessionFactoryBean.setDataSource(routingDataSource);
    //        sqlSessionFactoryBean.setMapperLocations(new PathMatchingResourcePatternResolver().getResources("classpath:mapper/*.xml"));
            return sqlSessionFactoryBean.getObject();
        }
    
        @Bean
        public PlatformTransactionManager platformTransactionManager() {
            return new DataSourceTransactionManager(routingDataSource);
        }
    }
    
    ```

    动态数据源

    ```java
    import org.springframework.jdbc.datasource.lookup.AbstractRoutingDataSource;
    import org.springframework.lang.Nullable;
    
    public class RoutingDataSource extends AbstractRoutingDataSource {
        @Nullable
        @Override
        protected Object determineCurrentLookupKey() {
            return DBContextHolder.get();
        }
    }
    
    ```

    数据源切换

    ```java
    
    import java.util.concurrent.atomic.AtomicInteger;
    
    public class DBContextHolder {
        private static final ThreadLocal<DBTypeEnum> contextHolder = new ThreadLocal<>();
    
        private static final AtomicInteger counter = new AtomicInteger(-1);
    
        public static void set(DBTypeEnum dbType) {
            contextHolder.set(dbType);
        }
    
        public static DBTypeEnum get() {
            return contextHolder.get();
        }
    
        public static void remove() {
            contextHolder.remove();
        }
    
        public static void master() {
            set(DBTypeEnum.MASTER);
            System.out.println("切换到master");
        }
    
        public static void slave() {
            //  轮询
            int index = counter.getAndIncrement() % 2;
            if (counter.get() > 9999) {
                counter.set(-1);
            }
            if (index == 0) {
                set(DBTypeEnum.SLAVE1);
                System.out.println("切换到slave1");
            } else {
                set(DBTypeEnum.SLAVE2);
                System.out.println("切换到slave2");
            }
        }
    
    }
    
    ```

    ## aop 切面

    

    ```java
    import org.aspectj.lang.annotation.After;
    import org.aspectj.lang.annotation.Aspect;
    import org.aspectj.lang.annotation.Before;
    import org.springframework.stereotype.Component;
    
    @Aspect
    @Component
    public class DataSourceAop {
        @Before("@annotation(com.example.demo.annotation.Master)")
        public void master() {
            DBContextHolder.master();
        }
    
        @Before("@annotation(com.example.demo.annotation.Slave)")
        public void slave() {
            DBContextHolder.slave();
        }
    
        @After("@annotation(com.example.demo.annotation.Slave)||@annotation(com.example.demo.annotation.Master)")
        public void afterSwitchDB() {
            DBContextHolder.remove();
        }
    }
    ```

    

    ## 使用方式

    ```java
    
    
    @Service
    public class UserServiceImpl implements UserService {
        @Autowired
        private UserMapper userMapper;
    
        /**
         * 主库写入数据
         * @param user
         * @return
         */
        @Override
        @Master
        public int save(User user) {
            return userMapper.save(user);
        }
    
        /**
         * 从库查询数据
         * @return
         */
        @Override
        @Slave
        public User get() {
            return userMapper.get();
        }
    }
    
    
    ```

    ```java
    
    public interface UserMapper {
    
        @Insert("insert into t_test(name)values(#{name})")
        int save(User user);
    
        @Select("select * from t_test where id>=(select floor(rand() * (select max(id) from t_test))) order by id limit 1")
        User get();
    }
    
    ```

    

    ## test

    ```java
    @SpringBootTest
    class DemoApplicationTests {
        @Autowired
        private UserService userService;
    
        @Test
        void testDb() throws InterruptedException {
            for (int i = 0; i < 10; i++) {
                User user = userService.get();
                System.out.println(user);
                TimeUnit.SECONDS.sleep(1);
                user.setName(UUID.randomUUID().toString());
                userService.save(user);
            }
        }
    
    }
    ```

    ```bash
    切换到slave2
    User{id=8, name='e87f057f-9dc9-4f59-8ff4-3a01682487d0'}
    切换到master
    切换到slave1
    User{id=2, name='d786167b-d29f-49eb-bb1d-dc48f01f761c'}
    切换到master
    切换到slave2
    User{id=6, name='c680a686-3ede-419d-bd66-de82966d5f96'}
    切换到master
    切换到slave1
    User{id=5, name='07d2a8d4-3414-49a0-bc37-065a1688a55f'}
    切换到master
    切换到slave2
    User{id=3, name='bac785d9-743d-48db-9b74-c2d51ba878fa'}
    切换到master
    切换到slave1
    User{id=7, name='0c3866a0-208b-4530-8c8d-9502dea753ff'}
    切换到master
    切换到slave2
    User{id=3, name='bac785d9-743d-48db-9b74-c2d51ba878fa'}
    切换到master
    切换到slave1
    User{id=8, name='e87f057f-9dc9-4f59-8ff4-3a01682487d0'}
    切换到master
    ```