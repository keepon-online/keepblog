var M=(p,y,n)=>new Promise((d,r)=>{var u=t=>{try{c(n.next(t))}catch(i){r(i)}},s=t=>{try{c(n.throw(t))}catch(i){r(i)}},c=t=>t.done?d(t.value):Promise.resolve(t.value).then(u,s);c((n=n.apply(p,y)).next())});import{d as P,r as _,L as f,Z as k,o as B,a1 as I,c as h,t as T,a as w,a6 as S,a3 as E,a7 as A,f as l,a4 as N,a5 as H,x as b,a9 as R,_ as G}from"./index-CRnhayG4.js";const K=["onClick","onMouseenter"],V={class:"cmd-icon"},q={class:"cmd-content"},z={class:"cmd-name"},O={class:"cmd-desc"},U={class:"cmd-badge"},W={key:0,class:"no-matches"},X=P({__name:"SlashCommand",props:{visible:{type:Boolean},position:{},filterText:{default:""}},emits:["select","close"],setup(p,{emit:y}){const n=p,d=y,r=_(null),u=_(null),s=_(0),c=[{id:"ai",label:"AI 智能续写",desc:"根据前文逻辑流式生成下一自然段",slash:"/ai",icon:"✨",action:"ai-continue"},{id:"outline",label:"AI 文章大纲",desc:"依据主题与受众定位智能生成多级技术大纲",slash:"/outline",icon:"📑",action:"ai-outline"},{id:"code",label:"代码块 (Code Block)",desc:"插入带语言高亮的代码围栏",slash:"/code",icon:"💻",template:"```go\n// 在此编写代码\n\n```\n"},{id:"table",label:"数据表格 (Table)",desc:"插入 3x3 标准 Markdown 数据表格",slash:"/table",icon:"📊",template:`| 参数名 | 类型 | 必填 | 说明 |
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
`}],t=f(()=>{const e=(n.filterText||"").trim().toLowerCase();return e?c.filter(a=>a.slash.toLowerCase().includes(e)||a.label.toLowerCase().includes(e)||a.id.toLowerCase().includes(e)):c}),i=f(()=>{var e;return((e=n.position)==null?void 0:e.placement)||"bottom"}),D=f(()=>{const e=i.value==="top";return{top:`${n.position.top}px`,left:`${n.position.left}px`,transform:e?"translateY(-100%)":"translateY(0)"}}),$=f(()=>{var m,v;if(typeof window=="undefined")return"280px";const e=window.innerHeight,a=i.value==="top",o=(v=(m=n.position)==null?void 0:m.top)!=null?v:300;if(a){const g=o-66;return`${Math.max(120,Math.min(280,g))}px`}else{const g=e-o-66;return`${Math.max(120,Math.min(280,g))}px`}});k(()=>n.filterText,()=>{s.value=0}),k(()=>n.visible,e=>{e&&(s.value=0)}),k(s,()=>M(null,null,function*(){if(yield R(),!u.value)return;const e=u.value.querySelector(".command-item.active");e&&e.scrollIntoView({block:"nearest"})}));const x=e=>{d("select",e)},C=e=>{if(n.visible)if(e.key==="ArrowDown")e.preventDefault(),e.stopPropagation(),t.value.length>0&&(s.value=(s.value+1)%t.value.length);else if(e.key==="ArrowUp")e.preventDefault(),e.stopPropagation(),t.value.length>0&&(s.value=(s.value-1+t.value.length)%t.value.length);else if(e.key==="Enter"){e.preventDefault(),e.stopPropagation();const a=t.value[s.value];a&&x(a)}else e.key==="Escape"&&(e.preventDefault(),e.stopPropagation(),d("close"))},L=e=>{n.visible&&r.value&&!r.value.contains(e.target)&&d("close")};return B(()=>{window.addEventListener("keydown",C,!0),window.addEventListener("pointerdown",L,!0)}),I(()=>{window.removeEventListener("keydown",C,!0),window.removeEventListener("pointerdown",L,!0)}),(e,a)=>p.visible?(w(),h("div",{key:0,ref_key:"menuRef",ref:r,class:A(["slash-command-menu",[`placement-${i.value}`]]),style:E(D.value),onClick:a[0]||(a[0]=S(()=>{},["stop"]))},[a[1]||(a[1]=l("div",{class:"menu-header"},[l("span",{class:"menu-title"},"⚡ 快速命令 (Slash Command)"),l("span",{class:"menu-tip"},"↑↓ 切换 · 回车确认 · Esc 退出")],-1)),l("div",{ref_key:"menuListRef",ref:u,class:"menu-list",style:E({maxHeight:$.value})},[(w(!0),h(N,null,H(t.value,(o,m)=>(w(),h("div",{key:o.id,class:A(["command-item",{active:s.value===m}]),onClick:v=>x(o),onMouseenter:v=>s.value=m},[l("span",V,b(o.icon),1),l("div",q,[l("span",z,b(o.label),1),l("span",O,b(o.desc),1)]),l("span",U,b(o.slash),1)],42,K))),128)),t.value.length===0?(w(),h("div",W," 未找到匹配命令 ")):T("",!0)],4)],6)):T("",!0)}}),Z=G(X,[["__scopeId","data-v-8b5b2a1b"]]);export{Z as default};
