package net

import (
	"bingo/pkg/utils"
	"fmt"
	"regexp"
	"strings"
)

var (
	// HTMLTags 标签	描述
	// <!--...-->	定义注释。
	// <!DOCTYPE> 	定义文档类型。
	// <a>	定义锚。
	// <abbr>	定义缩写。
	// <acronym>	定义只取首字母的缩写。
	// <address>	定义文档作者或拥有者的联系信息。
	// <applet>	不赞成使用。定义嵌入的 applet。
	// <area>	定义图像映射内部的区域。
	// <article>	定义文章。
	// <aside>	定义页面内容之外的内容。
	// <audio>	定义声音内容。
	// <b>	定义粗体字。
	// <base>	定义页面中所有链接的默认地址或默认目标。
	// <basefont>	不赞成使用。定义页面中文本的默认字体、颜色或尺寸。
	// <bdi>	定义文本的文本方向，使其脱离其周围文本的方向设置。
	// <bdo>	定义文字方向。
	// <big>	定义大号文本。
	// <blockquote>	定义长的引用。
	// <body>	定义文档的主体。
	// <br>	定义简单的折行。
	// <button>	定义按钮 (push button)。
	// <canvas>	定义图形。
	// <caption>	定义表格标题。
	// <center>	不赞成使用。定义居中文本。
	// <cite>	定义引用(citation)。
	// <code>	定义计算机代码文本。
	// <col>	定义表格中一个或多个列的属性值。
	// <colgroup>	定义表格中供格式化的列组。
	// <command>	定义命令按钮。
	// <datalist>	定义下拉列表。
	// <dd>	定义定义列表中项目的描述。
	// <del>	定义被删除文本。
	// <details>	定义元素的细节。
	// <dir>	不赞成使用。定义目录列表。
	// <div>	定义文档中的节。
	// <dfn>	定义定义项目。
	// <dialog>	定义对话框或窗口。
	// <dl>	定义定义列表。
	// <dt>	定义定义列表中的项目。
	// <em>	定义强调文本。
	// <embed>	定义外部交互内容或插件。
	// <fieldset>	定义围绕表单中元素的边框。
	// <figcaption>	定义 figure 元素的标题。
	// <figure>	定义媒介内容的分组，以及它们的标题。
	// <font>	不赞成使用。定义文字的字体、尺寸和颜色。
	// <footer>	定义 section 或 page 的页脚。
	// <form>	定义供用户输入的 HTML 表单。
	// <frame>	定义框架集的窗口或框架。
	// <frameset>	定义框架集。
	// <h1> to <h6>	定义 HTML 标题。
	// <head>	定义关于文档的信息。
	// <header>	定义 section 或 page 的页眉。
	// <hr>	定义水平线。
	// <html>	定义 HTML 文档。
	// <i>	定义斜体字。
	// <iframe>	定义内联框架。
	// <img>	定义图像。
	// <input>	定义输入控件。
	// <ins>	定义被插入文本。
	// <isindex>	不赞成使用。定义与文档相关的可搜索索引。
	// <kbd>	定义键盘文本。
	// <keygen>	定义生成密钥。
	// <label>	定义 input 元素的标注。
	// <legend>	定义 fieldset 元素的标题。
	// <li>	定义列表的项目。
	// <link>	定义文档与外部资源的关系。
	// <map>	定义图像映射。
	// <mark>	定义有记号的文本。
	// <menu>	定义命令的列表或菜单。
	// <menuitem>	定义用户可以从弹出菜单调用的命令/菜单项目。
	// <meta>	定义关于 HTML 文档的元信息。
	// <meter>	定义预定义范围内的度量。
	// <nav>	定义导航链接。
	// <noframes>	定义针对不支持框架的用户的替代内容。
	// <noscript>	定义针对不支持客户端脚本的用户的替代内容。
	// <object>	定义内嵌对象。
	// <ol>	定义有序列表。
	// <optgroup>	定义选择列表中相关选项的组合。
	// <option>	定义选择列表中的选项。
	// <output>	定义输出的一些类型。
	// <p>	定义段落。
	// <param>	定义对象的参数。
	// <pre>	定义预格式文本。
	// <progress>	定义任何类型的任务的进度。
	// <q>	定义短的引用。
	// <rp>	定义若浏览器不支持 ruby 元素显示的内容。
	// <rt>	定义 ruby 注释的解释。
	// <ruby>	定义 ruby 注释。
	// <s>	不赞成使用。定义加删除线的文本。
	// <samp>	定义计算机代码样本。
	// <script>	定义客户端脚本。
	// <section>	定义 section。
	// <select>	定义选择列表（下拉列表）。
	// <small>	定义小号文本。
	// <source>	定义媒介源。
	// <span>	定义文档中的节。
	// <strike>	不赞成使用。定义加删除线文本。
	// <strong>	定义强调文本。
	// <style>	定义文档的样式信息。
	// <sub>	定义下标文本。
	// <summary>	为 <details> 元素定义可见的标题。
	// <sup>	定义上标文本。
	// <table>	定义表格。
	// <tbody>	定义表格中的主体内容。
	// <td>	定义表格中的单元。
	// <textarea>	定义多行的文本输入控件。
	// <tfoot>	定义表格中的表注内容（脚注）。
	// <th>	定义表格中的表头单元格。
	// <thead>	定义表格中的表头内容。
	// <time>	定义日期/时间。
	// <title>	定义文档的标题。
	// <tr>	定义表格中的行。
	// <track>	定义用在媒体播放器中的文本轨道。
	// <tt>	定义打字机文本。
	// <u>	不赞成使用。定义下划线文本。
	// <ul>	定义无序列表。
	// <var>	定义文本的变量部分。
	// <video>	定义视频。
	// <wbr>	定义可能的换行符。
	// <xmp>	不赞成使用。定义预格式文本。

	HTMLTags = []string{
		"a",
		"abbr",
		"acronym",
		"address",
		"applet",
		"area",
		"article",
		"aside",
		"audio",
		"b",
		"base",
		"basefont",
		"bdi",
		"bdo",
		"big",
		"blockquote",
		"body",
		"br",
		"button",
		"canvas",
		"caption",
		"center",
		"cite",
		"code",
		"col",
		"colgroup",
		"command",
		"datalist",
		"dd",
		"del",
		"details",
		"dir",
		"div",
		"dfn",
		"dialog",
		"dl",
		"dt",
		"em",
		"embed",
		"fieldset",
		"figcaption",
		"figure",
		"font",
		"footer",
		"form",
		"frame",
		"frameset",
		"h1",
		"head",
		"header",
		"hr",
		"html",
		"i",
		"iframe",
		"img",
		"input",
		"ins",
		"isindex",
		"kbd",
		"keygen",
		"label",
		"legend",
		"li",
		"link",
		"map",
		"mark",
		"menu",
		"menuitem",
		"meta",
		"meter",
		"nav",
		"noframes",
		"noscript",
		"object",
		"ol",
		"optgroup",
		"option",
		"output",
		"p",
		"param",
		"pre",
		"progress",
		"q",
		"rp",
		"rt",
		"ruby",
		"s",
		"samp",
		"script",
		"section",
		"select",
		"small",
		"source",
		"span",
		"strike",
		"strong",
		"style",
		"sub",
		"summary",
		"sup",
		"table",
		"tbody",
		"td",
		"textarea",
		"tfoot",
		"th",
		"thead",
		"time",
		"title",
		"tr",
		"track",
		"tt",
		"u",
		"ul",
		"var",
		"video",
		"wbr",
		"xmp",
	}
	// HTMLString 显示结果	描述	实体名称	实体编号
	//  	空格	&nbsp;	&#160;
	// <	小于号	&lt;	&#60;
	// >	大于号	&gt;	&#62;
	// &	和号	&amp;	&#38;
	// "	引号	&quot;	&#34;
	// '	撇号 	&apos; (IE不支持)	&#39;
	// ￠	分（cent）	&cent;	&#162;
	// £	镑（pound）	&pound;	&#163;
	// ¥	元（yen）	&yen;	&#165;
	// €	欧元（euro）	&euro;	&#8364;
	// §	小节	&sect;	&#167;
	// ©	版权（copyright）	&copy;	&#169;
	// ®	注册商标	&reg;	&#174;
	// ™	商标	&trade;	&#8482;
	// ×	乘号	&times;	&#215;
	// ÷	除号	&divide;	&#247;
	HTMLString = map[string]string{
		"&nbsp;":   " ",
		"&lt;":     "<",
		"&gt;":     ">",
		"&amp;":    "&",
		"&quot;":   "\"",
		"&apos;":   "'",
		"&cent;":   "￠",
		"&pound;":  "£",
		"&yen;":    "¥",
		"&euro;":   "€",
		"&sect;":   "§",
		"&copy;":   "©",
		"&reg;":    "®",
		"&trade;":  "™",
		"&times;":  "×",
		"&divide;": "÷",
	}
	CodeString = map[string]string{
		"&#160;":  " ",
		"&#60;":   "<",
		"&#62;":   ">",
		"&#38;":   "&",
		"&#34;":   "\"",
		"&#39;":   "'",
		"&#162;":  "￠",
		"&#163;":  "£",
		"&#165;":  "¥",
		"&#8364;": "€",
		"&#167;":  "§",
		"&#169;":  "©",
		"&#174;":  "®",
		"&#8482;": "™",
		"&#215;":  "×",
		"&#247;":  "÷",
	}
)

func trimHTML(input string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src := re.ReplaceAllStringFunc(input, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

// HTML2Puretext returns pure text without html tags and flags
func HTML2Puretext(input string) string {
	tmp := utils.Unicode2UTF8(input)
	// return TrimHtml(tmp)

	// log.Notice("1111111     %s", tmp)
	tmp = strings.Replace(tmp, "<br/>", "\n", -1)
	tmp = strings.Replace(tmp, "<br />", "\n", -1)
	tmp = strings.Replace(tmp, "</p>", "\n", -1)
	tmp = strings.Replace(tmp, "&#38;", "&", -1)
	tmp = strings.Replace(tmp, "&amp;", "&", -1)

	for k, s := range HTMLString {
		tmp = strings.Replace(tmp, k, s, -1)
	}
	for k, s := range CodeString {
		tmp = strings.Replace(tmp, k, s, -1)
	}

	return strings.TrimSpace(removeHTML(tmp))
}

func findTag(input string) (int, string) {
	// fmt.Println(input)
	ileft := strings.Index(input, "<")
	if ileft != -1 {
		rest := input[ileft+1:]
		// fmt.Println(rest)
		for _, t := range HTMLTags {
			if strings.HasPrefix(rest, t+" ") || strings.HasPrefix(rest, "/"+t+">") || strings.HasPrefix(rest, t+">") {
				// fmt.Printf("-------TAG: %s; %d------------------\n", t, ileft)
				// break
				return ileft, t
			}
		}
		if ileft == strings.LastIndex(input, "<") {
			// fmt.Println("==========未找到11111============")
			return -1, ""
		}
		return findTag(rest)
	}
	// fmt.Println("==========未找到2222222============")
	return -1, ""
}

func removeHTML(input string) string {
	ileft, t := findTag(input)
	if ileft >= 0 {
		ileft = strings.Index(input, "<"+t)
		rest := input[ileft+1:]
		// fmt.Println("=========>>>>>>" + rest)
		iright := strings.Index(rest, ">")
		// fmt.Printf("%d-%d: %s\n", ileft, iright, t)
		if iright >= 0 {
			s := rest[iright+1:]
			// fmt.Println(">>>>>>" + s)
			if ileft > 0 {
				s = input[:ileft] + s
				// fmt.Println("<<<<<<" + input[:ileft])
			}
			s = strings.Replace(s, "</"+t+">", "", -1)
			s = strings.Replace(s, "<"+t+">", "", -1)
			// fmt.Println("------>>>>>>" + s)
			return removeHTML(s)
		}
	}
	return input
}

func iRemoveHTML(input string) string {
	ileft, t := findTag(input)
	if ileft >= 0 {
		input = strings.Replace(input, "</"+t+">", "", -1)
		iright := strings.Index(input, ">") // OK
		// fmt.Printf("%d-%d: %s\n", ileft, iright, t)
		if ileft < iright { //OK
			// if iright >= 0 {
			// s := rest[iright+1:]
			s := input[iright+1:] //OK
			if ileft > 0 {
				s = input[:ileft] + s
			}
			fmt.Println("------>>>>>>" + s)
			return iRemoveHTML(s)
		}
	}
	return input
}
