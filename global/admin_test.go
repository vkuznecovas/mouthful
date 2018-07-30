package global_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/global"
)

const input = `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>admin</title><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="mobile-web-app-capable" content="yes"><meta name="apple-mobile-web-app-capable" content="yes"><link rel="manifest" href="/manifest.json"><meta name="theme-color" content="#673ab8"><link rel="shortcut icon" href="/favicon.ico"><link href="/style.c10b5.css" rel="stylesheet"></head><body><div id="app"><header class="header__3QGkI"><h1>Mouthful Admin Panel</h1></header><div class="mouthful_container__1eO2q"><div class="mouthful_login__fadGC">No comments yet!</div></div></div><script defer="defer" src="/bundle.fb3a6.js"></script><script>window.fetch||document.write('<script src="/polyfills.44897.js"><\/script>')</script></body></html>`
const prefix = `/test`
const expectedOutput = `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>admin</title><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="mobile-web-app-capable" content="yes"><meta name="apple-mobile-web-app-capable" content="yes"><link rel="manifest" href="/test/manifest.json"><meta name="theme-color" content="#673ab8"><link rel="shortcut icon" href="/test/favicon.ico"><link href="/test/style.c10b5.css" rel="stylesheet"></head><body><div id="app"><header class="header__3QGkI"><h1>Mouthful Admin Panel</h1></header><div class="mouthful_container__1eO2q"><div class="mouthful_login__fadGC">No comments yet!</div></div></div><script defer="defer" src="/test/bundle.fb3a6.js"></script><script>window.fetch||document.write('<script src="/test/polyfills.44897.js"><\/script>')</script></body></html>`

const scriptInput = `!function(t){function e(n){if(r[n])return r[n].exports;var o=r[n]={i:n,l:!1,exports:{}};return t[n].call(o.exports,o,o.exports,e),o.l=!0,o.exports}var n=window.webpackJsonp;window.webpackJsonp=function(e,r,i){for(var u,a,c=0,l=[];c<e.length;c++)a=e[c],o[a]&&l.push(o[a][0]),o[a]=0;for(u in r)Object.prototype.hasOwnProperty.call(r,u)&&(t[u]=r[u]);for(n&&n(e,r,i);l.length;)l.shift()()};var r={},o={1:0};e.e=function(t){function n(){a.onerror=a.onload=null,clearTimeout(c);var e=o[t];0!==e&&(e&&e[1](new Error("Loading chunk "+t+" failed.")),o[t]=void 0)}var r=o[t];if(0===r)return new Promise(function(t){t()});if(r)return r[2];var i=new Promise(function(e,n){r=o[t]=[e,n]});r[2]=i;var u=document.getElementsByTagName("head")[0],a=document.createElement("script");a.type="text/javascript",a.charset="utf-8",a.async=!0,a.timeout=12e4,e.nc&&a.setAttribute("nonce",e.nc),a.src=e.p+""+({0:"route-panel"}[t]||t)+".chunk."+{0:"2aa25"}[t]+".js";var c=setTimeout(n,12e4);return a.onerror=a.onload=n,u.appendChild(a),i},e.m=t,e.c=r,e.d=function(t,n,r){e.o(t,n)||Object.defineProperty(t,n,{configurable:!1,enumerable:!0,get:r})},e.n=function(t){var n=t&&t.__esModule?function(){return t.default}:function(){return t};return e.d(n,"a",n),n},e.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},e.p="/",e.oe=function(t){throw console.error(t),t},e(e.s="pwNi")}({"/QC5":function(t,e,n){"use strict";function r(t,e){for(var n in e)t[n]=e[n];return t}function o(t,e,n){var r,o=/(?:\?([^#]*))?(#.*)?$/,i=t.match(o),u={};if(i&&i[1])for(var c=i[1].split("&"),l=0;l<c.length;l++){var p=c[l].split("=");u[decodeURIComponent(p[0])]=decodeURIComponent(p.slice(1).join("="))}t=a(t.replace(o,"")),e=a(e||"");for(var s=Math.max(t.length,e.length),f=0;f<s;f++)if(e[f]&&":"===e[f].charAt(0)){var d=e[f].replace(/(^\:|[+*?]+$)/g,""),h=(e[f].match(/[+*?]+$/)||x)[0]||"",_=~h.indexOf("+"),v=~h.indexOf("*"),m=t[f]||"";if(!m&&!v&&(h.indexOf("?")<0||_)){r=!1;break}if(u[d]=decodeURIComponent(m),_||v){u[d]=t.slice(f).map(decodeURIComponent).join("/");break}}else if(e[f]!==t[f]){r=!1;break}return(!0===n.default||!1!==r)&&u}function i(t,e){return t.rank<e.rank?1:t.rank>e.rank?-1:t.index-e.index}function u(t,e){return t.index=e,t.rank=p(t),t.attributes}function a(t){return t.replace(/(^\/+|\/+$)/g,"").split("/")}function c(t){return":"==t.charAt(0)?1+"*+?".indexOf(t.charAt(t.length-1))||4:5}function l(t){return a(t).map(c).join("")}function p(t){return t.attributes.default?0:l(t.attributes.path)}function s(t){return null!=t.__preactattr_||"undefined"!=typeof Symbol&&null!=t[Symbol.for("preactattr")]}function f(t,e){void 0===e&&(e="push"),O&&O[e]?O[e](t):"undefined"!=typeof history&&history[e+"State"]&&history[e+"State"](null,null,t)}function d(){var t;return t=O&&O.location?O.location:O&&O.getCurrentLocation?O.getCurrentLocation():"undefined"!=typeof location?location:N,""+(t.pathname||"")+(t.search||"")}function h(t,e){return void 0===e&&(e=!1),"string"!=typeof t&&t.url&&(e=t.replace,t=t.url),_(t)&&f(t,e?"replace":"push"),v(t)}function _(t){for(var e=k.length;e--;)if(k[e].canRoute(t))return!0;return!1}function v(t){for(var e=!1,n=0;n<k.length;n++)!0===k[n].routeTo(t)&&(e=!0);for(var r=j.length;r--;)j[r](t);return e}function m(t){if(t&&t.getAttribute){var e=t.getAttribute("href"),n=t.getAttribute("target");if(e&&e.match(/^\//g)&&(!n||n.match(/^_?self$/i)))return h(e)}}function y(t){if(0==t.button)return m(t.currentTarget||t.target||this),b(t)}function b(t){return t&&(t.stopImmediatePropagation&&t.stopImmediatePropagation(),t.stopPropagation&&t.stopPropagation(),t.preventDefault()),!1}function g(t){if(!(t.ctrlKey||t.metaKey||t.altKey||t.shiftKey||0!==t.button)){var e=t.target;do{if("A"===String(e.nodeName).toUpperCase()&&e.getAttribute("href")&&s(e)){if(e.hasAttribute("native"))return;if(m(e))return b(t)}}while(e=e.parentNode)}}function w(){U||("function"==typeof addEventListener&&(O||addEventListener("popstate",function(){v(d())}),addEventListener("click",g)),U=!0)}Object.defineProperty(e,"__esModule",{value:!0}),n.d(e,"subscribers",function(){return j}),n.d(e,"getCurrentUrl",function(){return d}),n.d(e,"route",function(){return h}),n.d(e,"Router",function(){return M}),n.d(e,"Route",function(){return L}),n.d(e,"Link",function(){return P});var C=n("KM04"),x=(n.n(C),{}),O=null,k=[],j=[],N={},U=!1,M=function(t){function e(e){t.call(this,e),e.history&&(O=e.history),this.state={url:e.url||d()},w()}return t&&(e.__proto__=t),e.prototype=Object.create(t&&t.prototype),e.prototype.constructor=e,e.prototype.shouldComponentUpdate=function(t){return!0!==t.static||(t.url!==this.props.url||t.onChange!==this.props.onChange)},e.prototype.canRoute=function(t){return this.getMatchingChildren(this.props.children,t,!1).length>0},e.prototype.routeTo=function(t){return this._didRoute=!1,this.setState({url:t}),this.updating?this.canRoute(t):(this.forceUpdate(),this._didRoute)},e.prototype.componentWillMount=function(){k.push(this),this.updating=!0},e.prototype.componentDidMount=function(){var t=this;O&&(this.unlisten=O.listen(function(e){t.routeTo(""+(e.pathname||"")+(e.search||""))})),this.updating=!1},e.prototype.componentWillUnmount=function(){"function"==typeof this.unlisten&&this.unlisten(),k.splice(k.indexOf(this),1)},e.prototype.componentWillUpdate=function(){this.updating=!0},e.prototype.componentDidUpdate=function(){this.updating=!1},e.prototype.getMatchingChildren=function(t,e,n){return t.filter(u).sort(i).map(function(t){var i=o(e,t.attributes.path,t.attributes);if(i){if(!1!==n){var u={url:e,matches:i};return r(u,i),delete u.ref,delete u.key,Object(C.cloneElement)(t,u)}return t}}).filter(Boolean)},e.prototype.render=function(t,e){var n=t.children,r=t.onChange,o=e.url,i=this.getMatchingChildren(n,o,!0),u=i[0]||null;this._didRoute=!!u;var a=this.previousUrl;return o!==a&&(this.previousUrl=o,"function"==typeof r&&r({router:this,url:o,previous:a,active:i,current:u})),u},e}(C.Component),P=function(t){return Object(C.h)("a",r({onClick:y},t))},L=function(t){return Object(C.h)(t.component,t)};M.subscribers=j,M.getCurrentUrl=d,M.route=h,M.Router=M,M.Route=L,M.Link=P,e.default=M},"7N8r":function(t,e,n){"use strict";e.__esModule=!0,e.default=function(t){function e(){var e=this;r.Component.call(this);var n=function(t){e.setState({child:t&&t.default||t})},o=t(n);o&&o.then&&o.then(n)}return e.prototype=new r.Component,e.prototype.constructor=e,e.prototype.render=function(t,e){return(0,r.h)(e.child,t)},e};var r=n("KM04")},JkW7:function(t,e,n){"use strict";function r(t,e){if(!t)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return!e||"object"!=typeof e&&"function"!=typeof e?t:e}function o(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Super expression must either be null or a function, not "+typeof e);t.prototype=Object.create(e&&e.prototype,{constructor:{value:t,enumerable:!1,writable:!0,configurable:!0}}),e&&(Object.setPrototypeOf?Object.setPrototypeOf(t,e):t.__proto__=e)}function i(t){n.e(0).then(function(){t(n("5Rz5"))}.bind(null,n)).catch(n.oe)}function u(t,e){if(!(t instanceof e))throw new TypeError("Cannot call a class as a function")}function a(t,e){if(!t)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return!e||"object"!=typeof e&&"function"!=typeof e?t:e}function c(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Super expression must either be null or a function, not "+typeof e);t.prototype=Object.create(e&&e.prototype,{constructor:{value:t,enumerable:!1,writable:!0,configurable:!0}}),e&&(Object.setPrototypeOf?Object.setPrototypeOf(t,e):t.__proto__=e)}Object.defineProperty(e,"__esModule",{value:!0});var l=(n("rq4c"),n("KM04")),p=(n("sw5u"),n("u3et")),s=n.n(p),f=Object(l.h)("h1",null,"Mouthful Admin Panel"),d=function(t){function e(){return r(this,t.apply(this,arguments))}return o(e,t),e.prototype.render=function(){return Object(l.h)("header",{class:s.a.header},f)},e}(l.Component),h=n("7N8r"),_=n.n(h),v=_()(i),m=Object(l.h)("div",{id:"app"},Object(l.h)(d,null),Object(l.h)(v,null)),y=function(t){function e(){var n,r,o;u(this,e);for(var i=arguments.length,c=Array(i),l=0;l<i;l++)c[l]=arguments[l];return n=r=a(this,t.call.apply(t,[this].concat(c))),r.handleRoute=function(t){r.currentUrl=t.url},o=n,a(r,o)}return c(e,t),e.prototype.render=function(){return m},e}(l.Component);e.default=y},KM04:function(t){!function(){"use strict";function e(){}function n(t,n){var r,o,i,u,a=R;for(u=arguments.length;u-- >2;)E.push(arguments[u]);for(n&&null!=n.children&&(E.length||E.push(n.children),delete n.children);E.length;)if((o=E.pop())&&void 0!==o.pop)for(u=o.length;u--;)E.push(o[u]);else"boolean"==typeof o&&(o=null),(i="function"!=typeof t)&&(null==o?o="":"number"==typeof o?o+="":"string"!=typeof o&&(i=!1)),i&&r?a[a.length-1]+=o:a===R?a=[o]:a.push(o),r=i;var c=new e;return c.nodeName=t,c.children=a,c.attributes=null==n?void 0:n,c.key=null==n?void 0:n.key,void 0!==L.vnode&&L.vnode(c),c}function r(t,e){for(var n in e)t[n]=e[n];return t}function o(t,e){return n(t.nodeName,r(r({},t.attributes),e),arguments.length>2?[].slice.call(arguments,2):t.children)}function i(t){!t.__d&&(t.__d=!0)&&1==A.push(t)&&(L.debounceRendering||T)(u)}function u(){var t,e=A;for(A=[];t=e.pop();)t.__d&&j(t)}function a(t,e,n){return"string"==typeof e||"number"==typeof e?void 0!==t.splitText:"string"==typeof e.nodeName?!t._componentConstructor&&c(t,e.nodeName):n||t._componentConstructor===e.nodeName}function c(t,e){return t.__n===e||t.nodeName.toLowerCase()===e.toLowerCase()}function l(t){var e=r({},t.attributes);e.children=t.children;var n=t.nodeName.defaultProps;if(void 0!==n)for(var o in n)void 0===e[o]&&(e[o]=n[o]);return e}function p(t,e){var n=e?document.createElementNS("http://www.w3.org/2000/svg",t):document.createElement(t);return n.__n=t,n}function s(t){var e=t.parentNode;e&&e.removeChild(t)}function f(t,e,n,r,o){if("className"===e&&(e="class"),"key"===e);else if("ref"===e)n&&n(null),r&&r(t);else if("class"!==e||o)if("style"===e){if(r&&"string"!=typeof r&&"string"!=typeof n||(t.style.cssText=r||""),r&&"object"==typeof r){if("string"!=typeof n)for(var i in n)i in r||(t.style[i]="");for(var i in r)t.style[i]="number"==typeof r[i]&&!1===S.test(i)?r[i]+"px":r[i]}}else if("dangerouslySetInnerHTML"===e)r&&(t.innerHTML=r.__html||"");else if("o"==e[0]&&"n"==e[1]){var u=e!==(e=e.replace(/Capture$/,""));e=e.toLowerCase().substring(2),r?n||t.addEventListener(e,h,u):t.removeEventListener(e,h,u),(t.__l||(t.__l={}))[e]=r}else if("list"!==e&&"type"!==e&&!o&&e in t)d(t,e,null==r?"":r),null!=r&&!1!==r||t.removeAttribute(e);else{var a=o&&e!==(e=e.replace(/^xlink\:?/,""));null==r||!1===r?a?t.removeAttributeNS("http://www.w3.org/1999/xlink",e.toLowerCase()):t.removeAttribute(e):"function"!=typeof r&&(a?t.setAttributeNS("http://www.w3.org/1999/xlink",e.toLowerCase(),r):t.setAttribute(e,r))}else t.className=r||""}function d(t,e,n){try{t[e]=n}catch(t){}}function h(t){return this.__l[t.type](L.event&&L.event(t)||t)}function _(){for(var t;t=W.pop();)L.afterMount&&L.afterMount(t),t.componentDidMount&&t.componentDidMount()}function v(t,e,n,r,o,i){I++||(K=null!=o&&void 0!==o.ownerSVGElement,D=null!=t&&!("__preactattr_"in t));var u=m(t,e,n,r,i);return o&&u.parentNode!==o&&o.appendChild(u),--I||(D=!1,i||_()),u}function m(t,e,n,r,o){var i=t,u=K;if(null!=e&&"boolean"!=typeof e||(e=""),"string"==typeof e||"number"==typeof e)return t&&void 0!==t.splitText&&t.parentNode&&(!t._component||o)?t.nodeValue!=e&&(t.nodeValue=e):(i=document.createTextNode(e),t&&(t.parentNode&&t.parentNode.replaceChild(i,t),b(t,!0))),i.__preactattr_=!0,i;var a=e.nodeName;if("function"==typeof a)return N(t,e,n,r);if(K="svg"===a||"foreignObject"!==a&&K,a+="",(!t||!c(t,a))&&(i=p(a,K),t)){for(;t.firstChild;)i.appendChild(t.firstChild);t.parentNode&&t.parentNode.replaceChild(i,t),b(t,!0)}var l=i.firstChild,s=i.__preactattr_,f=e.children;if(null==s){s=i.__preactattr_={};for(var d=i.attributes,h=d.length;h--;)s[d[h].name]=d[h].value}return!D&&f&&1===f.length&&"string"==typeof f[0]&&null!=l&&void 0!==l.splitText&&null==l.nextSibling?l.nodeValue!=f[0]&&(l.nodeValue=f[0]):(f&&f.length||null!=l)&&y(i,f,n,r,D||null!=s.dangerouslySetInnerHTML),w(i,e.attributes,s),K=u,i}function y(t,e,n,r,o){var i,u,c,l,p,f=t.childNodes,d=[],h={},_=0,v=0,y=f.length,g=0,w=e?e.length:0;if(0!==y)for(var C=0;C<y;C++){var x=f[C],O=x.__preactattr_,k=w&&O?x._component?x._component.__k:O.key:null;null!=k?(_++,h[k]=x):(O||(void 0!==x.splitText?!o||x.nodeValue.trim():o))&&(d[g++]=x)}if(0!==w)for(var C=0;C<w;C++){l=e[C],p=null;var k=l.key;if(null!=k)_&&void 0!==h[k]&&(p=h[k],h[k]=void 0,_--);else if(!p&&v<g)for(i=v;i<g;i++)if(void 0!==d[i]&&a(u=d[i],l,o)){p=u,d[i]=void 0,i===g-1&&g--,i===v&&v++;break}p=m(p,l,n,r),c=f[C],p&&p!==t&&p!==c&&(null==c?t.appendChild(p):p===c.nextSibling?s(c):t.insertBefore(p,c))}if(_)for(var C in h)void 0!==h[C]&&b(h[C],!1);for(;v<=g;)void 0!==(p=d[g--])&&b(p,!1)}function b(t,e){var n=t._component;n?U(n):(null!=t.__preactattr_&&t.__preactattr_.ref&&t.__preactattr_.ref(null),!1!==e&&null!=t.__preactattr_||s(t),g(t))}function g(t){for(t=t.lastChild;t;){var e=t.previousSibling;b(t,!0),t=e}}function w(t,e,n){var r;for(r in n)e&&null!=e[r]||null==n[r]||f(t,r,n[r],n[r]=void 0,K);for(r in e)"children"===r||"innerHTML"===r||r in n&&e[r]===("value"===r||"checked"===r?t[r]:n[r])||f(t,r,n[r],n[r]=e[r],K)}function C(t){var e=t.constructor.name;($[e]||($[e]=[])).push(t)}function x(t,e,n){var r,o=$[t.name];if(t.prototype&&t.prototype.render?(r=new t(e,n),M.call(r,e,n)):(r=new M(e,n),r.constructor=t,r.render=O),o)for(var i=o.length;i--;)if(o[i].constructor===t){r.__b=o[i].__b,o.splice(i,1);break}return r}function O(t,e,n){return this.constructor(t,n)}function k(t,e,n,r,o){t.__x||(t.__x=!0,(t.__r=e.ref)&&delete e.ref,(t.__k=e.key)&&delete e.key,!t.base||o?t.componentWillMount&&t.componentWillMount():t.componentWillReceiveProps&&t.componentWillReceiveProps(e,r),r&&r!==t.context&&(t.__c||(t.__c=t.context),t.context=r),t.__p||(t.__p=t.props),t.props=e,t.__x=!1,0!==n&&(1!==n&&!1===L.syncComponentUpdates&&t.base?i(t):j(t,1,o)),t.__r&&t.__r(t))}function j(t,e,n,o){if(!t.__x){var i,u,a,c=t.props,p=t.state,s=t.context,f=t.__p||c,d=t.__s||p,h=t.__c||s,m=t.base,y=t.__b,g=m||y,w=t._component,C=!1;if(m&&(t.props=f,t.state=d,t.context=h,2!==e&&t.shouldComponentUpdate&&!1===t.shouldComponentUpdate(c,p,s)?C=!0:t.componentWillUpdate&&t.componentWillUpdate(c,p,s),t.props=c,t.state=p,t.context=s),t.__p=t.__s=t.__c=t.__b=null,t.__d=!1,!C){i=t.render(c,p,s),t.getChildContext&&(s=r(r({},s),t.getChildContext()));var O,N,M=i&&i.nodeName;if("function"==typeof M){var P=l(i);u=w,u&&u.constructor===M&&P.key==u.__k?k(u,P,1,s,!1):(O=u,t._component=u=x(M,P,s),u.__b=u.__b||y,u.__u=t,k(u,P,0,s,!1),j(u,1,n,!0)),N=u.base}else a=g,O=w,O&&(a=t._component=null),(g||1===e)&&(a&&(a._component=null),N=v(a,i,s,n||!m,g&&g.parentNode,!0));if(g&&N!==g&&u!==w){var E=g.parentNode;E&&N!==E&&(E.replaceChild(N,g),O||(g._component=null,b(g,!1)))}if(O&&U(O),t.base=N,N&&!o){for(var R=t,T=t;T=T.__u;)(R=T).base=N;N._component=R,N._componentConstructor=R.constructor}}if(!m||n?W.unshift(t):C||(t.componentDidUpdate&&t.componentDidUpdate(f,d,h),L.afterUpdate&&L.afterUpdate(t)),null!=t.__h)for(;t.__h.length;)t.__h.pop().call(t);I||o||_()}}function N(t,e,n,r){for(var o=t&&t._component,i=o,u=t,a=o&&t._componentConstructor===e.nodeName,c=a,p=l(e);o&&!c&&(o=o.__u);)c=o.constructor===e.nodeName;return o&&c&&(!r||o._component)?(k(o,p,3,n,r),t=o.base):(i&&!a&&(U(i),t=u=null),o=x(e.nodeName,p,n),t&&!o.__b&&(o.__b=t,u=null),k(o,p,1,n,r),t=o.base,u&&t!==u&&(u._component=null,b(u,!1))),t}function U(t){L.beforeUnmount&&L.beforeUnmount(t);var e=t.base;t.__x=!0,t.componentWillUnmount&&t.componentWillUnmount(),t.base=null;var n=t._component;n?U(n):e&&(e.__preactattr_&&e.__preactattr_.ref&&e.__preactattr_.ref(null),t.__b=e,s(e),C(t),g(e)),t.__r&&t.__r(null)}function M(t,e){this.__d=!0,this.context=e,this.props=t,this.state=this.state||{}}function P(t,e,n){return v(n,t,{},!1,e,!1)}var L={},E=[],R=[],T="function"==typeof Promise?Promise.resolve().then.bind(Promise.resolve()):setTimeout,S=/acit|ex(?:s|g|n|p|$)|rph|ows|mnc|ntw|ine[ch]|zoo|^ord/i,A=[],W=[],I=0,K=!1,D=!1,$={};r(M.prototype,{setState:function(t,e){var n=this.state;this.__s||(this.__s=r({},n)),r(n,"function"==typeof t?t(n,this.props):t),e&&(this.__h=this.__h||[]).push(e),i(this)},forceUpdate:function(t){t&&(this.__h=this.__h||[]).push(t),j(this,2)},render:function(){}});var V={h:n,createElement:n,cloneElement:o,Component:M,render:P,rerender:u,options:L};t.exports=V}()},pwNi:function(t,e,n){"use strict";var r=n("KM04");"serviceWorker"in navigator&&"https:"===location.protocol&&navigator.serviceWorker.register(n.p+"sw.js");var o=function(t){return t&&t.default?t.default:t};if("function"==typeof o(n("JkW7"))){var i=document.body.firstElementChild,u=function(){var t=o(n("JkW7"));i=(0,r.render)((0,r.h)(t),document.body,i)};u()}},rq4c:function(){},sw5u:function(t,e,n){"use strict";function r(t,e){var n={};for(var r in t)e.indexOf(r)>=0||Object.prototype.hasOwnProperty.call(t,r)&&(n[r]=t[r]);return n}function o(t,e){if(!t)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return!e||"object"!=typeof e&&"function"!=typeof e?t:e}function i(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Super expression must either be null or a function, not "+typeof e);t.prototype=Object.create(e&&e.prototype,{constructor:{value:t,enumerable:!1,writable:!0,configurable:!0}}),e&&(Object.setPrototypeOf?Object.setPrototypeOf(t,e):t.__proto__=e)}Object.defineProperty(e,"__esModule",{value:!0}),e.Link=e.Match=void 0;var u=Object.assign||function(t){for(var e=1;e<arguments.length;e++){var n=arguments[e];for(var r in n)Object.prototype.hasOwnProperty.call(n,r)&&(t[r]=n[r])}return t},a=n("KM04"),c=n("/QC5"),l=e.Match=function(t){function e(){for(var e,n,r,i=arguments.length,u=Array(i),a=0;a<i;a++)u[a]=arguments[a];return e=n=o(this,t.call.apply(t,[this].concat(u))),n.update=function(t){n.nextUrl=t,n.setState({})},r=e,o(n,r)}return i(e,t),e.prototype.componentDidMount=function(){c.subscribers.push(this.update)},e.prototype.componentWillUnmount=function(){c.subscribers.splice(c.subscribers.indexOf(this.update)>>>0,1)},e.prototype.render=function(t){var e=this.nextUrl||(0,c.getCurrentUrl)(),n=e.replace(/\?.+$/,"");return this.nextUrl=null,t.children[0]&&t.children[0]({url:e,path:n,matches:n===t.path})},e}(a.Component),p=function(t){var e=t.activeClassName,n=t.path,o=r(t,["activeClassName","path"]);return(0,a.h)(l,{path:n||o.href},function(t){var n=t.matches;return(0,a.h)(c.Link,u({},o,{class:[o.class||o.className,n&&e].filter(Boolean).join(" ")}))})};e.Link=p,e.default=l,l.Link=p},u3et:function(t){t.exports={header:"header__3QGkI",active:"active__3gItZ"}}});
//# sourceMappingURL=bundle.b629d.js.map`

const expectedScriptOutput = `!function(t){function e(n){if(r[n])return r[n].exports;var o=r[n]={i:n,l:!1,exports:{}};return t[n].call(o.exports,o,o.exports,e),o.l=!0,o.exports}var n=window.webpackJsonp;window.webpackJsonp=function(e,r,i){for(var u,a,c=0,l=[];c<e.length;c++)a=e[c],o[a]&&l.push(o[a][0]),o[a]=0;for(u in r)Object.prototype.hasOwnProperty.call(r,u)&&(t[u]=r[u]);for(n&&n(e,r,i);l.length;)l.shift()()};var r={},o={1:0};e.e=function(t){function n(){a.onerror=a.onload=null,clearTimeout(c);var e=o[t];0!==e&&(e&&e[1](new Error("Loading chunk "+t+" failed.")),o[t]=void 0)}var r=o[t];if(0===r)return new Promise(function(t){t()});if(r)return r[2];var i=new Promise(function(e,n){r=o[t]=[e,n]});r[2]=i;var u=document.getElementsByTagName("head")[0],a=document.createElement("script");a.type="text/javascript",a.charset="utf-8",a.async=!0,a.timeout=12e4,e.nc&&a.setAttribute("nonce",e.nc),a.src=e.p+""+({0:"route-panel"}[t]||t)+".chunk."+{0:"2aa25"}[t]+".js";var c=setTimeout(n,12e4);return a.onerror=a.onload=n,u.appendChild(a),i},e.m=t,e.c=r,e.d=function(t,n,r){e.o(t,n)||Object.defineProperty(t,n,{configurable:!1,enumerable:!0,get:r})},e.n=function(t){var n=t&&t.__esModule?function(){return t.default}:function(){return t};return e.d(n,"a",n),n},e.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},e.p="/test/",e.oe=function(t){throw console.error(t),t},e(e.s="pwNi")}({"/QC5":function(t,e,n){"use strict";function r(t,e){for(var n in e)t[n]=e[n];return t}function o(t,e,n){var r,o=/(?:\?([^#]*))?(#.*)?$/,i=t.match(o),u={};if(i&&i[1])for(var c=i[1].split("&"),l=0;l<c.length;l++){var p=c[l].split("=");u[decodeURIComponent(p[0])]=decodeURIComponent(p.slice(1).join("="))}t=a(t.replace(o,"")),e=a(e||"");for(var s=Math.max(t.length,e.length),f=0;f<s;f++)if(e[f]&&":"===e[f].charAt(0)){var d=e[f].replace(/(^\:|[+*?]+$)/g,""),h=(e[f].match(/[+*?]+$/)||x)[0]||"",_=~h.indexOf("+"),v=~h.indexOf("*"),m=t[f]||"";if(!m&&!v&&(h.indexOf("?")<0||_)){r=!1;break}if(u[d]=decodeURIComponent(m),_||v){u[d]=t.slice(f).map(decodeURIComponent).join("/");break}}else if(e[f]!==t[f]){r=!1;break}return(!0===n.default||!1!==r)&&u}function i(t,e){return t.rank<e.rank?1:t.rank>e.rank?-1:t.index-e.index}function u(t,e){return t.index=e,t.rank=p(t),t.attributes}function a(t){return t.replace(/(^\/+|\/+$)/g,"").split("/")}function c(t){return":"==t.charAt(0)?1+"*+?".indexOf(t.charAt(t.length-1))||4:5}function l(t){return a(t).map(c).join("")}function p(t){return t.attributes.default?0:l(t.attributes.path)}function s(t){return null!=t.__preactattr_||"undefined"!=typeof Symbol&&null!=t[Symbol.for("preactattr")]}function f(t,e){void 0===e&&(e="push"),O&&O[e]?O[e](t):"undefined"!=typeof history&&history[e+"State"]&&history[e+"State"](null,null,t)}function d(){var t;return t=O&&O.location?O.location:O&&O.getCurrentLocation?O.getCurrentLocation():"undefined"!=typeof location?location:N,""+(t.pathname||"")+(t.search||"")}function h(t,e){return void 0===e&&(e=!1),"string"!=typeof t&&t.url&&(e=t.replace,t=t.url),_(t)&&f(t,e?"replace":"push"),v(t)}function _(t){for(var e=k.length;e--;)if(k[e].canRoute(t))return!0;return!1}function v(t){for(var e=!1,n=0;n<k.length;n++)!0===k[n].routeTo(t)&&(e=!0);for(var r=j.length;r--;)j[r](t);return e}function m(t){if(t&&t.getAttribute){var e=t.getAttribute("href"),n=t.getAttribute("target");if(e&&e.match(/^\//g)&&(!n||n.match(/^_?self$/i)))return h(e)}}function y(t){if(0==t.button)return m(t.currentTarget||t.target||this),b(t)}function b(t){return t&&(t.stopImmediatePropagation&&t.stopImmediatePropagation(),t.stopPropagation&&t.stopPropagation(),t.preventDefault()),!1}function g(t){if(!(t.ctrlKey||t.metaKey||t.altKey||t.shiftKey||0!==t.button)){var e=t.target;do{if("A"===String(e.nodeName).toUpperCase()&&e.getAttribute("href")&&s(e)){if(e.hasAttribute("native"))return;if(m(e))return b(t)}}while(e=e.parentNode)}}function w(){U||("function"==typeof addEventListener&&(O||addEventListener("popstate",function(){v(d())}),addEventListener("click",g)),U=!0)}Object.defineProperty(e,"__esModule",{value:!0}),n.d(e,"subscribers",function(){return j}),n.d(e,"getCurrentUrl",function(){return d}),n.d(e,"route",function(){return h}),n.d(e,"Router",function(){return M}),n.d(e,"Route",function(){return L}),n.d(e,"Link",function(){return P});var C=n("KM04"),x=(n.n(C),{}),O=null,k=[],j=[],N={},U=!1,M=function(t){function e(e){t.call(this,e),e.history&&(O=e.history),this.state={url:e.url||d()},w()}return t&&(e.__proto__=t),e.prototype=Object.create(t&&t.prototype),e.prototype.constructor=e,e.prototype.shouldComponentUpdate=function(t){return!0!==t.static||(t.url!==this.props.url||t.onChange!==this.props.onChange)},e.prototype.canRoute=function(t){return this.getMatchingChildren(this.props.children,t,!1).length>0},e.prototype.routeTo=function(t){return this._didRoute=!1,this.setState({url:t}),this.updating?this.canRoute(t):(this.forceUpdate(),this._didRoute)},e.prototype.componentWillMount=function(){k.push(this),this.updating=!0},e.prototype.componentDidMount=function(){var t=this;O&&(this.unlisten=O.listen(function(e){t.routeTo(""+(e.pathname||"")+(e.search||""))})),this.updating=!1},e.prototype.componentWillUnmount=function(){"function"==typeof this.unlisten&&this.unlisten(),k.splice(k.indexOf(this),1)},e.prototype.componentWillUpdate=function(){this.updating=!0},e.prototype.componentDidUpdate=function(){this.updating=!1},e.prototype.getMatchingChildren=function(t,e,n){return t.filter(u).sort(i).map(function(t){var i=o(e,t.attributes.path,t.attributes);if(i){if(!1!==n){var u={url:e,matches:i};return r(u,i),delete u.ref,delete u.key,Object(C.cloneElement)(t,u)}return t}}).filter(Boolean)},e.prototype.render=function(t,e){var n=t.children,r=t.onChange,o=e.url,i=this.getMatchingChildren(n,o,!0),u=i[0]||null;this._didRoute=!!u;var a=this.previousUrl;return o!==a&&(this.previousUrl=o,"function"==typeof r&&r({router:this,url:o,previous:a,active:i,current:u})),u},e}(C.Component),P=function(t){return Object(C.h)("a",r({onClick:y},t))},L=function(t){return Object(C.h)(t.component,t)};M.subscribers=j,M.getCurrentUrl=d,M.route=h,M.Router=M,M.Route=L,M.Link=P,e.default=M},"7N8r":function(t,e,n){"use strict";e.__esModule=!0,e.default=function(t){function e(){var e=this;r.Component.call(this);var n=function(t){e.setState({child:t&&t.default||t})},o=t(n);o&&o.then&&o.then(n)}return e.prototype=new r.Component,e.prototype.constructor=e,e.prototype.render=function(t,e){return(0,r.h)(e.child,t)},e};var r=n("KM04")},JkW7:function(t,e,n){"use strict";function r(t,e){if(!t)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return!e||"object"!=typeof e&&"function"!=typeof e?t:e}function o(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Super expression must either be null or a function, not "+typeof e);t.prototype=Object.create(e&&e.prototype,{constructor:{value:t,enumerable:!1,writable:!0,configurable:!0}}),e&&(Object.setPrototypeOf?Object.setPrototypeOf(t,e):t.__proto__=e)}function i(t){n.e(0).then(function(){t(n("5Rz5"))}.bind(null,n)).catch(n.oe)}function u(t,e){if(!(t instanceof e))throw new TypeError("Cannot call a class as a function")}function a(t,e){if(!t)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return!e||"object"!=typeof e&&"function"!=typeof e?t:e}function c(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Super expression must either be null or a function, not "+typeof e);t.prototype=Object.create(e&&e.prototype,{constructor:{value:t,enumerable:!1,writable:!0,configurable:!0}}),e&&(Object.setPrototypeOf?Object.setPrototypeOf(t,e):t.__proto__=e)}Object.defineProperty(e,"__esModule",{value:!0});var l=(n("rq4c"),n("KM04")),p=(n("sw5u"),n("u3et")),s=n.n(p),f=Object(l.h)("h1",null,"Mouthful Admin Panel"),d=function(t){function e(){return r(this,t.apply(this,arguments))}return o(e,t),e.prototype.render=function(){return Object(l.h)("header",{class:s.a.header},f)},e}(l.Component),h=n("7N8r"),_=n.n(h),v=_()(i),m=Object(l.h)("div",{id:"app"},Object(l.h)(d,null),Object(l.h)(v,null)),y=function(t){function e(){var n,r,o;u(this,e);for(var i=arguments.length,c=Array(i),l=0;l<i;l++)c[l]=arguments[l];return n=r=a(this,t.call.apply(t,[this].concat(c))),r.handleRoute=function(t){r.currentUrl=t.url},o=n,a(r,o)}return c(e,t),e.prototype.render=function(){return m},e}(l.Component);e.default=y},KM04:function(t){!function(){"use strict";function e(){}function n(t,n){var r,o,i,u,a=R;for(u=arguments.length;u-- >2;)E.push(arguments[u]);for(n&&null!=n.children&&(E.length||E.push(n.children),delete n.children);E.length;)if((o=E.pop())&&void 0!==o.pop)for(u=o.length;u--;)E.push(o[u]);else"boolean"==typeof o&&(o=null),(i="function"!=typeof t)&&(null==o?o="":"number"==typeof o?o+="":"string"!=typeof o&&(i=!1)),i&&r?a[a.length-1]+=o:a===R?a=[o]:a.push(o),r=i;var c=new e;return c.nodeName=t,c.children=a,c.attributes=null==n?void 0:n,c.key=null==n?void 0:n.key,void 0!==L.vnode&&L.vnode(c),c}function r(t,e){for(var n in e)t[n]=e[n];return t}function o(t,e){return n(t.nodeName,r(r({},t.attributes),e),arguments.length>2?[].slice.call(arguments,2):t.children)}function i(t){!t.__d&&(t.__d=!0)&&1==A.push(t)&&(L.debounceRendering||T)(u)}function u(){var t,e=A;for(A=[];t=e.pop();)t.__d&&j(t)}function a(t,e,n){return"string"==typeof e||"number"==typeof e?void 0!==t.splitText:"string"==typeof e.nodeName?!t._componentConstructor&&c(t,e.nodeName):n||t._componentConstructor===e.nodeName}function c(t,e){return t.__n===e||t.nodeName.toLowerCase()===e.toLowerCase()}function l(t){var e=r({},t.attributes);e.children=t.children;var n=t.nodeName.defaultProps;if(void 0!==n)for(var o in n)void 0===e[o]&&(e[o]=n[o]);return e}function p(t,e){var n=e?document.createElementNS("http://www.w3.org/2000/svg",t):document.createElement(t);return n.__n=t,n}function s(t){var e=t.parentNode;e&&e.removeChild(t)}function f(t,e,n,r,o){if("className"===e&&(e="class"),"key"===e);else if("ref"===e)n&&n(null),r&&r(t);else if("class"!==e||o)if("style"===e){if(r&&"string"!=typeof r&&"string"!=typeof n||(t.style.cssText=r||""),r&&"object"==typeof r){if("string"!=typeof n)for(var i in n)i in r||(t.style[i]="");for(var i in r)t.style[i]="number"==typeof r[i]&&!1===S.test(i)?r[i]+"px":r[i]}}else if("dangerouslySetInnerHTML"===e)r&&(t.innerHTML=r.__html||"");else if("o"==e[0]&&"n"==e[1]){var u=e!==(e=e.replace(/Capture$/,""));e=e.toLowerCase().substring(2),r?n||t.addEventListener(e,h,u):t.removeEventListener(e,h,u),(t.__l||(t.__l={}))[e]=r}else if("list"!==e&&"type"!==e&&!o&&e in t)d(t,e,null==r?"":r),null!=r&&!1!==r||t.removeAttribute(e);else{var a=o&&e!==(e=e.replace(/^xlink\:?/,""));null==r||!1===r?a?t.removeAttributeNS("http://www.w3.org/1999/xlink",e.toLowerCase()):t.removeAttribute(e):"function"!=typeof r&&(a?t.setAttributeNS("http://www.w3.org/1999/xlink",e.toLowerCase(),r):t.setAttribute(e,r))}else t.className=r||""}function d(t,e,n){try{t[e]=n}catch(t){}}function h(t){return this.__l[t.type](L.event&&L.event(t)||t)}function _(){for(var t;t=W.pop();)L.afterMount&&L.afterMount(t),t.componentDidMount&&t.componentDidMount()}function v(t,e,n,r,o,i){I++||(K=null!=o&&void 0!==o.ownerSVGElement,D=null!=t&&!("__preactattr_"in t));var u=m(t,e,n,r,i);return o&&u.parentNode!==o&&o.appendChild(u),--I||(D=!1,i||_()),u}function m(t,e,n,r,o){var i=t,u=K;if(null!=e&&"boolean"!=typeof e||(e=""),"string"==typeof e||"number"==typeof e)return t&&void 0!==t.splitText&&t.parentNode&&(!t._component||o)?t.nodeValue!=e&&(t.nodeValue=e):(i=document.createTextNode(e),t&&(t.parentNode&&t.parentNode.replaceChild(i,t),b(t,!0))),i.__preactattr_=!0,i;var a=e.nodeName;if("function"==typeof a)return N(t,e,n,r);if(K="svg"===a||"foreignObject"!==a&&K,a+="",(!t||!c(t,a))&&(i=p(a,K),t)){for(;t.firstChild;)i.appendChild(t.firstChild);t.parentNode&&t.parentNode.replaceChild(i,t),b(t,!0)}var l=i.firstChild,s=i.__preactattr_,f=e.children;if(null==s){s=i.__preactattr_={};for(var d=i.attributes,h=d.length;h--;)s[d[h].name]=d[h].value}return!D&&f&&1===f.length&&"string"==typeof f[0]&&null!=l&&void 0!==l.splitText&&null==l.nextSibling?l.nodeValue!=f[0]&&(l.nodeValue=f[0]):(f&&f.length||null!=l)&&y(i,f,n,r,D||null!=s.dangerouslySetInnerHTML),w(i,e.attributes,s),K=u,i}function y(t,e,n,r,o){var i,u,c,l,p,f=t.childNodes,d=[],h={},_=0,v=0,y=f.length,g=0,w=e?e.length:0;if(0!==y)for(var C=0;C<y;C++){var x=f[C],O=x.__preactattr_,k=w&&O?x._component?x._component.__k:O.key:null;null!=k?(_++,h[k]=x):(O||(void 0!==x.splitText?!o||x.nodeValue.trim():o))&&(d[g++]=x)}if(0!==w)for(var C=0;C<w;C++){l=e[C],p=null;var k=l.key;if(null!=k)_&&void 0!==h[k]&&(p=h[k],h[k]=void 0,_--);else if(!p&&v<g)for(i=v;i<g;i++)if(void 0!==d[i]&&a(u=d[i],l,o)){p=u,d[i]=void 0,i===g-1&&g--,i===v&&v++;break}p=m(p,l,n,r),c=f[C],p&&p!==t&&p!==c&&(null==c?t.appendChild(p):p===c.nextSibling?s(c):t.insertBefore(p,c))}if(_)for(var C in h)void 0!==h[C]&&b(h[C],!1);for(;v<=g;)void 0!==(p=d[g--])&&b(p,!1)}function b(t,e){var n=t._component;n?U(n):(null!=t.__preactattr_&&t.__preactattr_.ref&&t.__preactattr_.ref(null),!1!==e&&null!=t.__preactattr_||s(t),g(t))}function g(t){for(t=t.lastChild;t;){var e=t.previousSibling;b(t,!0),t=e}}function w(t,e,n){var r;for(r in n)e&&null!=e[r]||null==n[r]||f(t,r,n[r],n[r]=void 0,K);for(r in e)"children"===r||"innerHTML"===r||r in n&&e[r]===("value"===r||"checked"===r?t[r]:n[r])||f(t,r,n[r],n[r]=e[r],K)}function C(t){var e=t.constructor.name;($[e]||($[e]=[])).push(t)}function x(t,e,n){var r,o=$[t.name];if(t.prototype&&t.prototype.render?(r=new t(e,n),M.call(r,e,n)):(r=new M(e,n),r.constructor=t,r.render=O),o)for(var i=o.length;i--;)if(o[i].constructor===t){r.__b=o[i].__b,o.splice(i,1);break}return r}function O(t,e,n){return this.constructor(t,n)}function k(t,e,n,r,o){t.__x||(t.__x=!0,(t.__r=e.ref)&&delete e.ref,(t.__k=e.key)&&delete e.key,!t.base||o?t.componentWillMount&&t.componentWillMount():t.componentWillReceiveProps&&t.componentWillReceiveProps(e,r),r&&r!==t.context&&(t.__c||(t.__c=t.context),t.context=r),t.__p||(t.__p=t.props),t.props=e,t.__x=!1,0!==n&&(1!==n&&!1===L.syncComponentUpdates&&t.base?i(t):j(t,1,o)),t.__r&&t.__r(t))}function j(t,e,n,o){if(!t.__x){var i,u,a,c=t.props,p=t.state,s=t.context,f=t.__p||c,d=t.__s||p,h=t.__c||s,m=t.base,y=t.__b,g=m||y,w=t._component,C=!1;if(m&&(t.props=f,t.state=d,t.context=h,2!==e&&t.shouldComponentUpdate&&!1===t.shouldComponentUpdate(c,p,s)?C=!0:t.componentWillUpdate&&t.componentWillUpdate(c,p,s),t.props=c,t.state=p,t.context=s),t.__p=t.__s=t.__c=t.__b=null,t.__d=!1,!C){i=t.render(c,p,s),t.getChildContext&&(s=r(r({},s),t.getChildContext()));var O,N,M=i&&i.nodeName;if("function"==typeof M){var P=l(i);u=w,u&&u.constructor===M&&P.key==u.__k?k(u,P,1,s,!1):(O=u,t._component=u=x(M,P,s),u.__b=u.__b||y,u.__u=t,k(u,P,0,s,!1),j(u,1,n,!0)),N=u.base}else a=g,O=w,O&&(a=t._component=null),(g||1===e)&&(a&&(a._component=null),N=v(a,i,s,n||!m,g&&g.parentNode,!0));if(g&&N!==g&&u!==w){var E=g.parentNode;E&&N!==E&&(E.replaceChild(N,g),O||(g._component=null,b(g,!1)))}if(O&&U(O),t.base=N,N&&!o){for(var R=t,T=t;T=T.__u;)(R=T).base=N;N._component=R,N._componentConstructor=R.constructor}}if(!m||n?W.unshift(t):C||(t.componentDidUpdate&&t.componentDidUpdate(f,d,h),L.afterUpdate&&L.afterUpdate(t)),null!=t.__h)for(;t.__h.length;)t.__h.pop().call(t);I||o||_()}}function N(t,e,n,r){for(var o=t&&t._component,i=o,u=t,a=o&&t._componentConstructor===e.nodeName,c=a,p=l(e);o&&!c&&(o=o.__u);)c=o.constructor===e.nodeName;return o&&c&&(!r||o._component)?(k(o,p,3,n,r),t=o.base):(i&&!a&&(U(i),t=u=null),o=x(e.nodeName,p,n),t&&!o.__b&&(o.__b=t,u=null),k(o,p,1,n,r),t=o.base,u&&t!==u&&(u._component=null,b(u,!1))),t}function U(t){L.beforeUnmount&&L.beforeUnmount(t);var e=t.base;t.__x=!0,t.componentWillUnmount&&t.componentWillUnmount(),t.base=null;var n=t._component;n?U(n):e&&(e.__preactattr_&&e.__preactattr_.ref&&e.__preactattr_.ref(null),t.__b=e,s(e),C(t),g(e)),t.__r&&t.__r(null)}function M(t,e){this.__d=!0,this.context=e,this.props=t,this.state=this.state||{}}function P(t,e,n){return v(n,t,{},!1,e,!1)}var L={},E=[],R=[],T="function"==typeof Promise?Promise.resolve().then.bind(Promise.resolve()):setTimeout,S=/acit|ex(?:s|g|n|p|$)|rph|ows|mnc|ntw|ine[ch]|zoo|^ord/i,A=[],W=[],I=0,K=!1,D=!1,$={};r(M.prototype,{setState:function(t,e){var n=this.state;this.__s||(this.__s=r({},n)),r(n,"function"==typeof t?t(n,this.props):t),e&&(this.__h=this.__h||[]).push(e),i(this)},forceUpdate:function(t){t&&(this.__h=this.__h||[]).push(t),j(this,2)},render:function(){}});var V={h:n,createElement:n,cloneElement:o,Component:M,render:P,rerender:u,options:L};t.exports=V}()},pwNi:function(t,e,n){"use strict";var r=n("KM04");"serviceWorker"in navigator&&"https:"===location.protocol&&navigator.serviceWorker.register(n.p+"sw.js");var o=function(t){return t&&t.default?t.default:t};if("function"==typeof o(n("JkW7"))){var i=document.body.firstElementChild,u=function(){var t=o(n("JkW7"));i=(0,r.render)((0,r.h)(t),document.body,i)};u()}},rq4c:function(){},sw5u:function(t,e,n){"use strict";function r(t,e){var n={};for(var r in t)e.indexOf(r)>=0||Object.prototype.hasOwnProperty.call(t,r)&&(n[r]=t[r]);return n}function o(t,e){if(!t)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return!e||"object"!=typeof e&&"function"!=typeof e?t:e}function i(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Super expression must either be null or a function, not "+typeof e);t.prototype=Object.create(e&&e.prototype,{constructor:{value:t,enumerable:!1,writable:!0,configurable:!0}}),e&&(Object.setPrototypeOf?Object.setPrototypeOf(t,e):t.__proto__=e)}Object.defineProperty(e,"__esModule",{value:!0}),e.Link=e.Match=void 0;var u=Object.assign||function(t){for(var e=1;e<arguments.length;e++){var n=arguments[e];for(var r in n)Object.prototype.hasOwnProperty.call(n,r)&&(t[r]=n[r])}return t},a=n("KM04"),c=n("/QC5"),l=e.Match=function(t){function e(){for(var e,n,r,i=arguments.length,u=Array(i),a=0;a<i;a++)u[a]=arguments[a];return e=n=o(this,t.call.apply(t,[this].concat(u))),n.update=function(t){n.nextUrl=t,n.setState({})},r=e,o(n,r)}return i(e,t),e.prototype.componentDidMount=function(){c.subscribers.push(this.update)},e.prototype.componentWillUnmount=function(){c.subscribers.splice(c.subscribers.indexOf(this.update)>>>0,1)},e.prototype.render=function(t){var e=this.nextUrl||(0,c.getCurrentUrl)(),n=e.replace(/\?.+$/,"");return this.nextUrl=null,t.children[0]&&t.children[0]({url:e,path:n,matches:n===t.path})},e}(a.Component),p=function(t){var e=t.activeClassName,n=t.path,o=r(t,["activeClassName","path"]);return(0,a.h)(l,{path:n||o.href},function(t){var n=t.matches;return(0,a.h)(c.Link,u({},o,{class:[o.class||o.className,n&&e].filter(Boolean).join(" ")}))})};e.Link=p,e.default=l,l.Link=p},u3et:function(t){t.exports={header:"header__3QGkI",active:"active__3gItZ"}}});
//# sourceMappingURL=bundle.b629d.js.map`

func TestOverrideScriptRootInAdminHtml(t *testing.T) {
	filepath := "./t.html"

	var _, err = os.Stat(filepath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)
		defer file.Close()
		assert.Nil(t, err)

		_, err = file.WriteString(input)
		assert.Nil(t, err)

		file.Close()
	}

	err = global.OverrideScriptRootInAdminHTML("/test", filepath)
	assert.Nil(t, err)

	b, err := ioutil.ReadFile(filepath)
	assert.Nil(t, err)

	newHtml := string(b)
	assert.Equal(t, expectedOutput, newHtml)

	err = global.OverrideScriptRootInAdminHTML("/test", filepath)
	assert.Nil(t, err)

	b, err = ioutil.ReadFile(filepath)
	assert.Nil(t, err)

	newHtml = string(b)
	assert.Equal(t, expectedOutput, newHtml)

	err = os.Remove(filepath)
	assert.Nil(t, err)
}

func TestOverrideScriptPathInBundle(t *testing.T) {
	filepath := "./t.js"

	var _, err = os.Stat(filepath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)
		defer file.Close()
		assert.Nil(t, err)

		_, err = file.WriteString(expectedScriptOutput)
		assert.Nil(t, err)

		file.Close()
	}

	err = global.OverrideScriptPathInBundle("/test", filepath)
	assert.Nil(t, err)

	b, err := ioutil.ReadFile(filepath)
	assert.Nil(t, err)

	newHtml := string(b)
	assert.Equal(t, expectedScriptOutput, newHtml)

	err = global.OverrideScriptPathInBundle("/test", filepath)
	assert.Nil(t, err)

	b, err = ioutil.ReadFile(filepath)
	assert.Nil(t, err)

	newHtml = string(b)
	assert.Equal(t, expectedScriptOutput, newHtml)

	err = global.OverrideScriptPathInBundle("/", filepath)
	assert.Nil(t, err)

	b, err = ioutil.ReadFile(filepath)
	assert.Nil(t, err)

	newHtml = string(b)
	assert.Equal(t, scriptInput, newHtml)

	err = os.Remove(filepath)
	assert.Nil(t, err)
}

func TestOverrideScriptPathInBundle_Errors_On_No_File(t *testing.T) {
	err := global.OverrideScriptPathInBundle("/", "this does not exist")
	assert.Error(t, err)
}

func TestOverrideScriptPathInAdminHTML_Errors_On_No_File(t *testing.T) {
	err := global.OverrideScriptRootInAdminHTML("/", "this does not exist")
	assert.Error(t, err)
}

func TestOverrideScriptPathInBundleReturnsErrorOnWrongSliceLength(t *testing.T) {
	filepath := "./t.js"

	var _, err = os.Stat(filepath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)
		defer file.Close()
		assert.Nil(t, err)

		_, err = file.WriteString("Some string here")
		assert.Nil(t, err)

		file.Close()
	}
	err = global.OverrideScriptPathInBundle("/", filepath)
	assert.NotNil(t, err)
	assert.Equal(t, global.ErrCouldNotOverrideBundlePath, err)
	err = os.Remove(filepath)
	assert.Nil(t, err)
}
func TestFindAdminPanelChunkFilename(t *testing.T) {
	name := "bundle.2aa25.js"
	filepath := "./" + name
	var _, err = os.Stat(filepath)
	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)
		defer file.Close()
		assert.Nil(t, err)
	}
	parsedName, err := global.FindAdminPanelChunkFilename("./")
	assert.Nil(t, err)
	assert.Equal(t, name, parsedName)
	err = os.Remove(filepath)
	assert.Nil(t, err)
}

func TestRewriteAdminPanelScripts_CouldNotOverrideRoot(t *testing.T) {
	err := global.RewriteAdminPanelScripts("filepath")
	assert.NotNil(t, err)
	assert.Equal(t, "Couldn't override the static admin html root. Please check if the file index file(./static/index.html) is available and the mouthful user has the permissions to access it", err.Error())
}

func TestRewriteAdminPanelScripts_CouldNotFindBundle(t *testing.T) {
	filepath := "./static/index.html"

	var _, err = os.Stat(filepath)
	// create file if not exists
	if os.IsNotExist(err) {
		err := os.Mkdir("./static", 0777)
		assert.Nil(t, err)

		assert.Nil(t, err)
		file, err := os.Create(filepath)
		assert.Nil(t, err)
		defer file.Close()
		defer os.RemoveAll(filepath)
		defer os.RemoveAll("./static")
		_, err = file.WriteString(input)
		assert.Nil(t, err)

		file.Close()
	}
	err = global.RewriteAdminPanelScripts("/test")
	assert.NotNil(t, err)
	assert.Equal(t, "Couldn't find the admin panel chunk file. Please check if the budle file is available and the mouthful user has permissions to access it in the directory ./static", err.Error())
}

func TestRewriteAdminPanelScripts_CouldNotOverrideBundle(t *testing.T) {
	filepath := "./static/index.html"
	bundlePath := "./static/bundle.2aa25.js"
	var _, err = os.Stat(filepath)
	// create file if not exists
	if os.IsNotExist(err) {
		err := os.Mkdir("./static", 0777)
		assert.Nil(t, err)
		file, err := os.Create(bundlePath)
		defer os.RemoveAll(bundlePath)
		assert.Nil(t, err)
		file, err = os.Create(filepath)
		assert.Nil(t, err)
		defer file.Close()
		defer os.RemoveAll(filepath)
		defer os.RemoveAll("./static")
		_, err = file.WriteString(input)
		assert.Nil(t, err)
		file.Close()
	}
	err = global.RewriteAdminPanelScripts("/test")
	assert.NotNil(t, err)
	assert.Equal(t, "Couldn't override the static admin script path. Please check if the script file(./static/bundle.2aa25.js) is available and the mouthful user has the permissions to access it", err.Error())
}

func TestRewriteAdminPanelScripts_OK(t *testing.T) {
	filepath := "./static/index.html"
	bundlePath := "./static/bundle.2aa25.js"
	var _, err = os.Stat(filepath)
	// create file if not exists
	if os.IsNotExist(err) {
		err := os.Mkdir("./static", 0777)
		assert.Nil(t, err)
		bundleFile, err := os.Create(bundlePath)
		defer os.RemoveAll(bundlePath)
		assert.Nil(t, err)
		_, err = bundleFile.WriteString(scriptInput)
		assert.Nil(t, err)
		bundleFile.Close()
		file, err := os.Create(filepath)
		assert.Nil(t, err)
		defer file.Close()
		defer os.RemoveAll(filepath)
		defer os.RemoveAll("./static")
		_, err = file.WriteString(input)
		assert.Nil(t, err)
		file.Close()
	}
	err = global.RewriteAdminPanelScripts("/test")
	assert.Nil(t, err)
}

func TestFindAdminPanelChunkFilename_Errors_On_No_file(t *testing.T) {
	parsedName, err := global.FindAdminPanelChunkFilename("./")
	assert.Equal(t, "", parsedName)
	assert.Error(t, err)
}
