import{d as g,r as _,I as k,V as f,o as y,Y as x,c,$ as b,a as d,a4 as A,a1 as D,f as s,a2 as T,a3 as E,a5 as M,a6 as r,_ as $}from"./index-CYD25nR2.js";const B={class:"menu-list"},I=["onClick","onMouseenter"],L={class:"cmd-icon"},N={class:"cmd-content"},P={class:"cmd-name"},S={class:"cmd-desc"},G={class:"cmd-badge"},H={key:0,class:"no-matches"},K=g({__name:"SlashCommand",props:{visible:{type:Boolean},position:{},filterText:{}},emits:["select","close"],setup(l,{emit:w}){const i=l,u=w,n=_(0),m=[{id:"ai",label:"AI 智能续写",desc:"根据前文逻辑流式生成下一自然段",slash:"/ai",icon:"✨",action:"ai-continue"},{id:"outline",label:"AI 文章大纲",desc:"依据主题与受众定位智能生成多级技术大纲",slash:"/outline",icon:"📑",action:"ai-outline"},{id:"code",label:"代码块 (Code Block)",desc:"插入带语言高亮的代码围栏",slash:"/code",icon:"💻",template:"```go\n// 在此编写代码\n\n```\n"},{id:"table",label:"数据表格 (Table)",desc:"插入 3x3 标准 Markdown 数据表格",slash:"/table",icon:"📊",template:`| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :---: | :--- |
| name | string | 是 | 资源名称 |
| count | int | 否 | 数量限制 |
`},{id:"mermaid",label:"Mermaid 流程图",desc:"插入 Mermaid 图表架构图模板",slash:"/flow",icon:"🗺️",template:`\`\`\`mermaid
flowchart TD
    A[客户端请求] --> B[API 网关]
    B --> C{认证校验}
    C -->|通过| D[微服务集群]
    C -->|拒绝| E[返回 401]
\`\`\`
`},{id:"tip",label:"提示块 (Tip Alert)",desc:"插入 GitHub 风格绿色最佳实践提示块",slash:"/tip",icon:"💡",template:`> [!TIP]
> 这是一个高亮提示块，用于强调关键配置与最佳实践。

`},{id:"warning",label:"警告块 (Warning Alert)",desc:"插入 GitHub 风格黄色/橙色避坑警示",slash:"/warn",icon:"⚠️",template:`> [!WARNING]
> 生产环境操作需注意并发竞争与死锁风险。

`},{id:"note",label:"说明块 (Note Alert)",desc:"插入 GitHub 风格蓝色背景补充说明",slash:"/note",icon:"📝",template:`> [!NOTE]
> 本章节仅适用于 Linux 内核 5.4 及以上版本。

`},{id:"math",label:"数学公式块 (KaTeX)",desc:"插入独立 KaTeX 块级数学公式",slash:"/math",icon:"🔢",template:`$$
E = mc^2
$$
`}],a=k(()=>{const e=(i.filterText||"").trim().toLowerCase();return e?m.filter(t=>t.slash.toLowerCase().includes(e)||t.label.toLowerCase().includes(e)||t.id.toLowerCase().includes(e)):m});f(()=>i.filterText,()=>{n.value=0}),f(()=>i.visible,e=>{e&&(n.value=0)});const p=e=>{u("select",e)},v=e=>{if(i.visible)if(e.key==="ArrowDown")e.preventDefault(),e.stopPropagation(),a.value.length>0&&(n.value=(n.value+1)%a.value.length);else if(e.key==="ArrowUp")e.preventDefault(),e.stopPropagation(),a.value.length>0&&(n.value=(n.value-1+a.value.length)%a.value.length);else if(e.key==="Enter"){e.preventDefault(),e.stopPropagation();const t=a.value[n.value];t&&p(t)}else e.key==="Escape"&&(e.preventDefault(),e.stopPropagation(),u("close"))};return y(()=>{window.addEventListener("keydown",v,!0)}),x(()=>{window.removeEventListener("keydown",v,!0)}),(e,t)=>l.visible?(d(),c("div",{key:0,class:"slash-command-menu",style:D({top:`${l.position.top}px`,left:`${l.position.left}px`}),onClick:t[0]||(t[0]=A(()=>{},["stop"]))},[t[1]||(t[1]=s("div",{class:"menu-header"},[s("span",{class:"menu-title"},"⚡ 快速命令 (Slash Command)"),s("span",{class:"menu-tip"},"↑↓ 切换 · 回车确认 · Esc 退出")],-1)),s("div",B,[(d(!0),c(T,null,E(a.value,(o,h)=>(d(),c("div",{key:o.id,class:M(["command-item",{active:n.value===h}]),onClick:C=>p(o),onMouseenter:C=>n.value=h},[s("span",L,r(o.icon),1),s("div",N,[s("span",P,r(o.label),1),s("span",S,r(o.desc),1)]),s("span",G,r(o.slash),1)],42,I))),128)),a.value.length===0?(d(),c("div",H," 未找到匹配命令 ")):b("",!0)])],4)):b("",!0)}}),z=$(K,[["__scopeId","data-v-75252b27"]]);export{z as default};
