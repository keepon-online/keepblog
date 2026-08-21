package category

import (
	"testing"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/testutil"
)

func TestGetCategories_WithPostCount(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewCategoryService(testutil.NewTestDB(t))

	cats, err := s.GetCategories()
	if err != nil {
		t.Fatalf("GetCategories 报错: %v", err)
	}
	// 无文章的分类（隐藏分类）不出现在统计里；cat1=10 篇，cat2=2 篇
	if len(cats) != 2 {
		t.Fatalf("len = %d, want 2: %+v", len(cats), cats)
	}
	byName := map[string]uint32{}
	for _, c := range cats {
		byName[c.CategoryName] = c.Total
	}
	if byName["Go"] != 10 {
		t.Errorf("Go 文章数 = %d, want 10", byName["Go"])
	}
	if byName["MySQL"] != 2 {
		t.Errorf("MySQL 文章数 = %d, want 2", byName["MySQL"])
	}
}

func TestGetCategoryList_ReturnsAll(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewCategoryService(testutil.NewTestDB(t))

	cats, err := s.GetCategoryList()
	if err != nil {
		t.Fatalf("GetCategoryList 报错: %v", err)
	}
	if len(cats) != 3 {
		t.Errorf("后台分类列表 len = %d, want 3（含停用分类）", len(cats))
	}
}

func TestGetCategory_ById(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewCategoryService(testutil.NewTestDB(t))

	cat, err := s.GetCategory(2)
	if err != nil {
		t.Fatalf("GetCategory(2) 报错: %v", err)
	}
	if cat.CategoryName != "MySQL" {
		t.Errorf("CategoryName = %q, want MySQL", cat.CategoryName)
	}
	if _, err := s.GetCategory(999); err == nil {
		t.Error("不存在的分类应报错")
	}
}

func TestCategoryCRUD(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewCategoryService(testutil.NewTestDB(t))

	created := model.Category{CategoryName: "测试新建分类", Note: "n", State: 1}
	if err := s.Save(created); err != nil {
		t.Fatalf("Save 报错: %v", err)
	}
	list, _ := s.GetCategoryList()
	var id uint32
	for _, c := range list {
		if c.CategoryName == "测试新建分类" {
			id = c.CategoryId
		}
	}
	if id == 0 {
		t.Fatal("新建分类未找到")
	}

	created.CategoryId = id
	created.Note = "updated"
	if err := s.Update(created); err != nil {
		t.Fatalf("Update 报错: %v", err)
	}
	got, _ := s.GetCategory(id)
	if got.Note != "updated" {
		t.Errorf("Note = %q, want updated", got.Note)
	}

	if err := s.Delete(int(id)); err != nil {
		t.Fatalf("Delete 报错: %v", err)
	}
	if _, err := s.GetCategory(id); err == nil {
		t.Error("删除后应查不到")
	}
}
