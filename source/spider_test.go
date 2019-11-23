package source

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/apex/log"
	"github.com/phpgao/proxy_pool/model"
	"strings"
	"testing"
	"time"
)

func init() {
	logger.Level = log.DebugLevel
}

func TestGetSpiders(t *testing.T) {
	//for _, c := range ListOfSpider {
	//	testSpiderFetch(c)
	//}
	testSpider := &site_digger{}
	testSpiderFetch(testSpider)
}

func testSpiderFetch(c Crawler) {
	newProxyChan := make(chan *model.HttpProxy, 100)
	c.SetProxyChan(newProxyChan)
	c.Run()
	timeout := time.After(30 * time.Second)
	for {
		select {
		case proxy := <-newProxyChan:
			fmt.Println(proxy)
			go func(p *model.HttpProxy) {
				flag := p.SimpleTcpTest(5)
				fmt.Println(flag)
			}(proxy)
		case <-timeout:
			fmt.Println("There's no more time to this. Exiting!")
			return
		}
	}
}

func TestXpath(t *testing.T) {
	html := `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns:wb="http://open.weibo.com/wb">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>鲲鹏免费HTTP代理列表（每3小时更新一次）_鲲鹏Web数据抓取 - 专业Web数据采集服务提供商</title>
<meta name="keywords" content="HTTP代理,S5代理,SOCKS5代理,免费代理,共享代理, 网页数据抓取,网页数据采集,网站数据抓取,网站数据采集,匿名HTTP代理提供商,破解反采集策略,验证码识别,手机App抓取,微信小程序数据抓取" />
<meta name="description" content="鲲鹏数据目前提供四种类型的代理解决方案:1）免费代理;2）美国共享代理;3）香港独享代理;4）国内代理.后面三种均为稳定高匿代理，是鲲鹏数据自己在机房搭建的代理，有专人维护，稳定性好" />
<link href="/css/main.css" rel="stylesheet" type="text/css" />
<link type="text/css" rel="stylesheet" href="/js/highlighter/css/SyntaxHighlighter.css"></link>
<script type="text/javascript" src="/js/aes.js"></script>
<script type="text/javascript" src="/js/pad-zeropadding.js"></script>
<script language="javascript">
var baidu_union_id = "4ee9a5640e0a11ea";
function decrypt(d){var b=CryptoJS.enc.Latin1.parse(baidu_union_id);var c=b;var a=CryptoJS.AES.decrypt(d,b,{iv:c,padding:CryptoJS.pad.ZeroPadding});return a.toString(CryptoJS.enc.Utf8)};
</script>
</head>
<div id="header">

	<div class="main">

		<img alt="鲲鹏Web数据抓取 - 专业Web数据采集服务提供商" src="/images/logo.png"/>

		<h2><a href="/">鲲鹏Web数据抓取 - 专业Web数据采集服务提供商</a></h2>

		<h3>选择我们，所有数据都是你的！</h3>
        <h4>服务热线：029-87553281</h4>

		<ul>

			<li id="index_nav"><a href="/"><span>首页</span></a></li>

			

      		<li><a href='/html/services/' title="服务项目"><span>服务项目</span></a></li>

      		

      		<li><a href='/html/feepolicy/' title="收费政策"><span>收费政策</span></a></li>

      		

      		<li><a href='/html/submit/' title="提交项目"><span>提交项目</span></a></li>

      		<li class='current'><a href='/html/articles/'><span>技术文章</span></a></li>

      		<li><a href='/html/contact/' title="联系我们"><span>联系我们</span></a></li>

      		

		</ul>

	</div>

</div>



<div id="primary">
	<div class="main">
	  <div class="leftbox">
	     <div class="smallbox">
	<h3><span style="float:right; font-size:11px; font-weight:normal;"><a href="/html/about.html">更多>></a></span>关于我们</h3>
	<p><img src="/images/inf_ico.gif"> 西安鲲之鹏网络信息技术有限公司从2010年开始专注于Web（网站）数据抓取领域。致力于为广大中国客户提供准确、快捷的数据采集相关服务。我们采用分布式系统架构，日采集网页数千万。我们拥有海量稳定高匿HTTP代理IP地址池，可以有效绕过各种反采集策略。</p>
    <p><img src="/images/services_ico.gif"> 您只需告诉我们您想抓取的网站是什么，您感兴趣的字段有哪些，你需要的数据是哪种格式，我们将为您做所有的工作，最后把数据（或程序）交付给你。</p>
<p><img src="/images/db_ico.gif"> 数据的格式可以是CSV、JSON、XML、ACCESS、SQLITE、MSSQL、MYSQL等等。</p>
</div>

<div class="smallbox">
	<h3>快捷导航</h3>
	<ul class="scnav">
        <li><a href="/html/faq.html" title="常见问题答疑">FAQ</a></li>
		<li><a href="/html/services/" title="服务项目介绍">服务项目</a></li>
		<li><a href="/html/part-cases.html" title="典型案例/项目实例">典型案例</a></li>
		<li><a href="/html/sample.html" title="示例数据">示例数据</a></li>
        <li><a href="http://www.data-shop.net" title="中国最专业数据分享、销售平台">数据超市</a></li>
		<li><a href="/html/articles/" title="技术文章/经验分享">技术文章</a></li>
	</ul>
	<div style="clear:both"></div>
</div>

<div class="indexarclist">
	<h3><span style="float:right; font-size:11px; font-weight:normal;"><a href="/html/articles/">更多>></a></span>技术文章</h3>
	<ul>
	<li> <a href="/html/articles/20191018/751.html" title="一例ssl pinning突破过程记录">一例ssl pinning突破过程记录</a></li>
<li> <a href="/html/articles/20190928/750.html" title="Facebook消息自动发送辅助工具演示">Facebook消息自动发送辅助工具演示</a></li>
<li> <a href="/html/articles/20190916/749.html" title="手机淘宝APP关键词搜索采集方案介绍">手机淘宝APP关键词搜索采集方案介绍</a></li>
<li> <a href="/html/articles/20190804/748.html" title="今日头条App广告采集器的实现">今日头条App广告采集器的实现</a></li>
<li> <a href="/html/articles/20190622/732.html" title="一例APK脱壳反编译寻找AES密钥过程记录">一例APK脱壳反编译寻找AES密钥过程记录</a></li>

	</ul>
 	<div style="clear:both"></div>
</div>
        
<div class="smallbox contactbox">
	<h3><span style="float:right; font-size:11px; font-weight:normal;"><a href="http://www.site-digger.com/html/weibo/" target="_blank">更多>></a></span>官方微博</h3>
	<div class="weibo_head">
		<div class="weibo_pic"><a href="http://weibo.com/kunzhipeng" target="_blank"></a></div>
		<div class="weibo_detail">
			<dl class="weibo_dev_name">
				<dt><a href="https://weibo.com/kunzhipeng" target="_blank" title="西安鲲之鹏官微">西安鲲之鹏</a><span class="weibo_ico" title="新浪认证"></span></dt>
				<dd>陕西 西安</dd>
			</dl>
			<p class="dev_follow">
				<span class="weibo_btn">
					<span>
						<em class="btn_ok"></em>
						<em><a href="https://weibo.com/kunzhipeng" target="_blank">加关注</a></em>
					</span>
				</span>
			</p>
		</div>
		<div class="clear"></div>
	</div>
	<div id="weibo_side">
		<ul class="weibo_sidebar list_side">
			<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						【经验分享】本文介绍了一例APP使用了非常规ssl pinning导致Fidder抓包失败，最终通过Frida HOOK成功解决，附源码。 &gt;&gt;&gt; <a href="http://www.site-digger.com/html/articles/20191018/751.html" title="http://www.site-digger.com/html/articles/20191018/751.html" target="_blank">http://www.site-digger.com/html/articles/20191018/751.html</a> ​​​​<br/><ul class="weibo_images"><li><a href="/uploads/weibo_images/mw690/0065K5msly1g82a8mu14zj30rs0im405.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g82a8mu14zj30rs0im405.jpg"/></a></li> <li><a href="/uploads/weibo_images/mw690/0065K5msly1g82a8mv33xj30ii0liabl.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g82a8mv33xj30ii0liabl.jpg"/></a></li> <li><a href="/uploads/weibo_images/mw690/0065K5msly1g82a8mto6yj30ip0h1gmo.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g82a8mto6yj30ip0h1gmo.jpg"/></a></li> <li><a href="/uploads/weibo_images/mw690/0065K5msly1g82a8mrii8j30gn06x3ye.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g82a8mrii8j30gn06x3ye.jpg"/></a></li> <li><a href="/uploads/weibo_images/mw690/0065K5msly1g82a8mwbccj318w0qgwl3.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g82a8mwbccj318w0qgwl3.jpg"/></a></li><div class="clear"></div></ul>
					</div>
					<h5 class="time">发布时间：2019-10-18 13:23:42</h5>
				</div>
				<div class="clear"></div>
			</li>
<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						 【经验分享】今日拨号服务器上某PPPOE拨号持续失败，经查日志：“pppoe: send (sendPacket): Network is down”，ip link 查看对应的虚拟网卡状态是DOWN，无法设置为UP（sudo ip link set dev v051802057684 up失败）。但同一个账号在另外一个机器上测试正常，怀疑可能是MAC地址的问题（例如冲突了），果断删掉虚拟网卡（ sudo sudo ip link del v051802057684），然后重建并指定一个不同的MAC，拨号成功！<br/><ul class="weibo_images"><li><a href="/uploads/weibo_images/mw690/0065K5msly1g7cq9d7tysj30m30gb0v7.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g7cq9d7tysj30m30gb0v7.jpg"/></a></li><div class="clear"></div></ul>
					</div>
					<h5 class="time">发布时间：2019-09-26 10:52:18</h5>
				</div>
				<div class="clear"></div>
			</li>
<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						 【经验分享】昨天下午办公室断了下电，服务器重启后，adb devices显示10台设备都是“?????? no permissions”，第一次遇到这种情况。重启服务器和移动设备问题依据。后来在askubuntu上看到有人提到试一下sudo adb devices，竟然立马识别了（<a href="https://askubuntu.com/questions/908306/adb-no-permissions-on-ubuntu-17-04" title="https://askubuntu.com/questions/908306/adb-no-permissions-on-ubuntu-17-04" target="_blank">https://askubuntu.com/questions/908306/adb-no-permissions-on-ubuntu-17-04</a>）。很诡异，之前用普通权限都一直正常着，为什么突然就没有权限了？
					</div>
					<h5 class="time">发布时间：2019-09-26 08:45:09</h5>
				</div>
				<div class="clear"></div>
			</li>
<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						【经验分享】Termux自带的am命令版本太低，竟然不支持force-stop（如图1示），经查得知pm clear命令也可以停止一个APP，而且还会将APP的用户数据清除掉（回到刚安装的状态），试了一下果然有效，运行之后&quot;/data/data/包名&quot;目录下只剩下lib目录了。注意：需要root权限。 &gt;&gt;&gt; <a href="https://stackoverflow.com/questions/3117095/stopping-an-android-app-from-console" title="https://stackoverflow.com/questions/3117095/stopping-an-android-app-from-console" target="_blank">https://stackoverflow.com/questions/3117095/stopping-an-android-app-from-console</a> ​​​​<br/><ul class="weibo_images"><li><a href="/uploads/weibo_images/mw690/0065K5msly1g74pc5av1dj30dc0jm3zl.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g74pc5av1dj30dc0jm3zl.jpg"/></a></li> <li><a href="/uploads/weibo_images/mw690/0065K5msly1g74phsts9uj30bx071dft.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g74phsts9uj30bx071dft.jpg"/></a></li><div class="clear"></div></ul>
					</div>
					<h5 class="time">发布时间：2019-09-19 12:20:13</h5>
				</div>
				<div class="clear"></div>
			</li>
<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						 【经验分享】&quot;adb shell 命令&quot;，如何让“命令”以root权限执行？<br>例如，某设备已root，但当执行adb shell rm /data/data/com.xxxx/cache时提示Permission denied。<br><br>解决方法：<br>adb shell &quot;su -c '[your command goes here]'&quot;<br>例如：<br>adb shell &quot;su -c 'rm /data/data/com.xxxx/cache'&quot;<br><br>参考文章&quot;Is there a way for me to run Adb shell as root without typing in 'su'?&quot; &gt;&gt;&gt; <a href="https://android.stackexchange.com/questions/5884/is-there-a-way-for-me-to-run-adb-shell-as-root-without-typing-in-su" title="https://android.stackexchange.com/questions/5884/is-there-a-way-for-me-to-run-adb-shell-as-root-without-typing-in-su" target="_blank">https://android.stackexchange.com/questions/5884/is-there-a-way-for-me-to-run-adb-shell-as-root-without-typing-in-su</a>
					</div>
					<h5 class="time">发布时间：2019-09-18 09:16:52</h5>
				</div>
				<div class="clear"></div>
			</li>
<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						【经验分享】关于小红书搜索结果数据抓取的限制：<br>1. 小红书微信小程序版之前是前600条可见，最近已被限制为前60条可见。<br>2. 小红书安卓APP版本限制为搜索结果前1000条可见。 ​​​​<br/><ul class="weibo_images"><li><a href="/uploads/weibo_images/mw690/0065K5msly1g6l9khmqd8j30vx0sw10y.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g6l9khmqd8j30vx0sw10y.jpg"/></a></li><div class="clear"></div></ul>
					</div>
					<h5 class="time">发布时间：2019-09-02 16:42:37</h5>
				</div>
				<div class="clear"></div>
			</li>
<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						【经验分享】通过ADB启动手机淘宝APP搜索，打开指定关键词的搜索结果列表，如下示例，打开搜索“HUAWEI”的结果：<br>adb shell am start -n com.taobao.taobao/com.taobao.search.SearchListActivity -d &quot;taobao://s.taobao.com/search?q=HUAWEI&quot; ​​​​<br/><ul class="weibo_images"><li><a href="/uploads/weibo_images/mw690/0065K5msly1g68ipb8ukyj30vf0mtwos.jpg" target="_blank"><img src="/uploads/weibo_images/thumb150/0065K5msly1g68ipb8ukyj30vf0mtwos.jpg"/></a></li> <li><a href="/uploads/weibo_images/mw690/0065K5msly1g68ipcx2jsj30td07zgos.jpg" target="_blank"><img src="/uploads/weibo_images/thumb150/0065K5msly1g68ipcx2jsj30td07zgos.jpg"/></a></li><div class="clear"></div></ul>
					</div>
					<h5 class="time">发布时间：2019-08-22 16:07:00</h5>
				</div>
				<div class="clear"></div>
			</li>
<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						【经验分享】手机淘宝APP商品搜索结果采集最新方案20190821(免登录)<br>（1）模拟操作APP操作，无需登录，不存在封账号；<br>（2）IP限制弱；<br>详见下附演示视频。&nbsp;&nbsp;<a  suda-uatrack="key=tblog_card&value=click_title:4407784365671052:1034-video:1034%3A4407783542426901:page_100606_manage:5581662372:4407784365671052:5581662372" title="西安鲲之鹏的微博视频" href="http://t.cn/AiQKo3lN" alt="http://t.cn/AiQKo3lN" action-type="feed_list_url" target="_blank" rel="noopener noreferrer"><i class="W_ficon ficon_cd_video">L</i>西安鲲之鹏的微博视频</a> ​​​​
					</div>
					<h5 class="time">发布时间：2019-08-21 17:52:28</h5>
				</div>
				<div class="clear"></div>
			</li>
<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						【经验分享】终于搞定了支付宝口碑App数据采集，有图有真相。 ​​​​<br/><ul class="weibo_images"><li><a href="/uploads/weibo_images/mw690/0065K5msly1g5tkoba9tuj30tl1frn9s.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g5tkoba9tuj30tl1frn9s.jpg"/></a></li><div class="clear"></div></ul>
					</div>
					<h5 class="time">发布时间：2019-08-09 17:52:21</h5>
				</div>
				<div class="clear"></div>
			</li>
<li>
				<div class="weibo_side_text">
					<div class="weibo_content">
						【经验分享】如何采集支付宝小程序的数据？adb模拟操作 + xposed Hook &quot;com.alipay.mobile.nebula.util.H5Utils.parseObject&quot; <br>如附图所示，成功获取服务端应答的JSON数据。 ​​​​<br/><ul class="weibo_images"><li><a href="/uploads/weibo_images/mw690/0065K5msly1g5r6jkzm35j30sp0vnaf2.jpg" target="_blank"><img src="/uploads/weibo_images/orj360/0065K5msly1g5r6jkzm35j30sp0vnaf2.jpg"/></a></li><div class="clear"></div></ul>
					</div>
					<h5 class="time">发布时间：2019-08-07 16:13:23</h5>
				</div>
				<div class="clear"></div>
			</li>
		
		</ul>
		<div class="more"><a href="http://www.site-digger.com/html/weibo/" target="_blank">查看更多</a></div>
	</div>
</div>

	  </div>
	  <div class="rightbox">
	  	<div id="content">
		  <div class="pos"> <span>当前位置:</span> <a href='/'>首页</a> > <a href='/html/articles/'>技术文章</a> >  </div>
		  <div id="content_detail">
			<div class="arctitle">鲲鹏免费HTTP代理列表（每3小时更新一次）</div>
			   <div class="arcinfo">
					发布时间：2011-05-16
			   </div>
			   <div class="arccontent">
					<div class="bdsharebuttonbox share"><a href="#" class="bds_more" data-cmd="more"></a><a href="#" class="bds_qzone" data-cmd="qzone"></a><a href="#" class="bds_tqq" data-cmd="tqq"></a><a href="#" class="bds_tsina" data-cmd="tsina"></a><a href="#" class="bds_weixin" data-cmd="weixin"></a></div>
                    <script>
                        window._bd_share_config={
                            "common":{
                                "bdSnsKey":{},
                                "bdText":"",
                                "bdMini":"2",
                                "bdPic":"",
                                "bdStyle":"0",
                                "bdSize":"16"
                            },
                            "share":{},
                            "image":{
                                "viewList":["qzone","tsina","tqq","renren","weixin"],
                                "viewText":"分享到：",
                                "viewSize":"16"
                            },
                            "selectShare":{
                                "bdContainerClass":null,
                                "bdSelectMiniList":["qzone","tsina","tqq","renren","weixin"]
                            }
                        };
                        with(document)0[(getElementsByTagName('head')[0]||body).appendChild(createElement('script')).src='http://bdimg.share.baidu.com/static/api/js/share.js?v=89860593.js?cdnversion='+~(-new Date()/36e5)];
                    </script> 
					<hr>
					<p style="color:red">如果您觉得免费的代理无法满足需求（比如稳定性、速度、匿名性等因素），鲲鹏数据还提供稳定、高匿、收费的HTTP代理方案，详情请查看这里 <a href="http://www.site-digger.com/html/proxies.html">http://www.site-digger.com/html/proxies.html</a></p>
<p><strong>来源：</strong>由鲲鹏数据爬虫程序采集于互联网。</p>
<p><strong>更新频率：</strong>每3小时更新一次。</p>
<p><strong>稳定性：</strong>不确定。&nbsp;</p>
<p><strong>实时状态：</strong>当前可用代理共有<font color="red" id="total_proxies_num">231</font>个，其中匿名代理<font color="red" id="anonymous_proxies_num">79</font>个。</p>
<p><style type="text/css">
#proxies_table{
	width:500px;
	border-collapse:collapse;
}
#proxies_table caption{
	font-weight:bold;
	margin-bottom:6px;
}
#proxies_table thead{
	font-weight: bold;
	background-color:#4791FF;
	color:#FFFFFF;
	text-align:center;
}
#proxies_table td, #proxies_table th{
	padding:5px;
	border:1px solid #CCCCCC;
}
</style></p>
<table id="proxies_table">
    <caption>上次更新时间: 2019-11-23 22:00:01，以下列表免费共享，仅1000条可见</caption>
    <thead>
        <tr>
            <th>代理</th>
            <th>类型</th>
            <th>国家</th>
            <th>延迟</th>
        </tr>
    </thead>
    <tbody><tr>
                    <td><script>document.write(decrypt("t8JiBA12jTxLLiqD3XnQXp8vyQlR1HV0VoWM5r9Vnqw="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("uDYtvW+VsuDexvtF6+IhCQpirbmH6/teW0okpamQ094="));</script></td>
                    <td>Transparent</td>
                    <td>加拿大</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("zkYC1FywqM+TEzQkfgFc0kRhou59JSNOehdBDIY8Rqw="));</script></td>
                    <td>Transparent</td>
                    <td>尼日利亚</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("vJCW95l2+ZQ+lRt7sfdCET4YDGT+Afc4zfi1hZhmGCQ="));</script></td>
                    <td>Transparent</td>
                    <td>香港</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("YghZkLybB1gRtIaWtD8PWlaXXHEuyhEXA35w5rGFnn0="));</script></td>
                    <td>Transparent</td>
                    <td>尼日利亚</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Yjvbdj0WS+JLpbQuVcN1LMXZW5k0dEw/aJkU9V0J2u0="));</script></td>
                    <td>Transparent</td>
                    <td>尼日利亚</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("8KLixk3dbniyf35+aQo8+g=="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("FhtKRHvri/ERqaAo6naHAB0EevW6QGNcPg9j5aHS5gE="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("BFGbSeTYMlq6QdTTAdekAjLQSv/9inP/dptTCcwQrCU="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("e1HPzsY2lvOAdCY1zRJgx8Dcbtq95NJ3nlHj0v+xA40="));</script></td>
                    <td>Transparent</td>
                    <td>巴西</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Wsg7vVDXhEydrLrocC54gl2KzRlpcyt3e8pJCXMxjsY="));</script></td>
                    <td>Anonymous</td>
                    <td>德国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("ePcfxOhGrUpJNwVPUGyTH/rbU06XLmniM7GnSZ+ni6o="));</script></td>
                    <td>Anonymous</td>
                    <td>德国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("kr0+0uq3J2YjKhHJu6gSRh0xkKsNLBQijaWDVQbWDI8="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("7eTFlScXSOZLiagxpM+aKQukYenkCZREa6r5Lv7kQRQ="));</script></td>
                    <td>Transparent</td>
                    <td>阿根廷</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("qnzSa9EtA2F5sHkW8b3gWAJfGL+mLaVqem4ckCyrl2w="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("bt/dehvpNlTyruGP56Vr1w=="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("V6gYQLK2CKK00EAV9xoy9hnzjyOO5o9PXzUZiUQdKYw="));</script></td>
                    <td>Anonymous</td>
                    <td>巴拿马</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("jMv5sWVVTSDtRXeWpBfeV9I8R2BAjHhOAMKHWMEc9Ok="));</script></td>
                    <td>Anonymous</td>
                    <td>香港</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("G8S8v5QmkRV60Z93l6j9B1mxfOXwLpfabafTYGfhqdc="));</script></td>
                    <td>Anonymous</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("cvi09KWB3I4aMnMbpX2qZw=="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("8vtPr3cgSxak1LJsIS/Rs7z4CFsJBzITEbXa3MRR7UA="));</script></td>
                    <td>Transparent</td>
                    <td>俄罗斯</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("jGUI8LzDQHU1pOwgw5MHtZ2oG2qMNNMcoAIS56l9Zqc="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>11</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("eC11MrlEOr6kty/nQi5BYyZhZ7UTOO32lIizGlxnwhs="));</script></td>
                    <td>Transparent</td>
                    <td>菲律宾</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("VMvP31FW/b0rspoFo3M63CHBbiFpThiXPvfFJph3j1Y="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("wm3HLpJx0yGMumSCEOCrDoQpDxZh1FrRJN6sTzcLYs8="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("1RyZ1WERw4BsPHq8rPGAN6gcVxtoW2doAgggYW84n90="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("U+whkUzy/bpuk0sKS9BjDyRocrih9ntX/ndmlpOjc4k="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>9</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("maMmDhHY1Jh2WSVZyDOT9XPARbdcRulltqfQrbnB77A="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("E89rNHPDVFgmSjY7GK7/GHKiNlYmuUmyDvaznaZqR1M="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("38Tz+jjpmKOO2HSCawQwVW6dh6UY1AHo+knWrHFmigs="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("8X65R9UOjhm16A0P4uiARdHUURoYArNP8rm6Ao94JUE="));</script></td>
                    <td>Anonymous</td>
                    <td>德国</th>
                    <th>12</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("HWFRm4wR1wC4052k0NLSvwjbi6+GHdxSL0QB8xRuxdk="));</script></td>
                    <td>Anonymous</td>
                    <td>印度</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("0bHHOyfNMo9SWpgSAAOoQ7N8n2SB/ihkN1BcKI/Gbks="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("qlqAf6e+rGndAstUcZ+jrmHzjaF6HRfrNCRdY6MW5Eo="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("GDbarw73Bowpe6wEt2FuZWprlT2BvpGaZ48cdqwCL1M="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("s3TVTgPWvBproSnESwb1AvMJEUQ/97IgW0syEnFhU+Q="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("W3YFP5LTdzcq8IQ2CThGSN8twYuMMfGjY9NwugAf41s="));</script></td>
                    <td>Transparent</td>
                    <td>阿根廷</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("LtsZsS127cmTtyVA/KxsDvb70oDZsA4jeSaa68ZaYDY="));</script></td>
                    <td>Transparent</td>
                    <td>墨西哥</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("idb3jTAiklEwsyNi7i58GzqX/wniDbtsFJ7DIULfR/4="));</script></td>
                    <td>Anonymous</td>
                    <td>加拿大</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("cRg1WBUXZ7j3341a37txH4NQj1AoduQQWTQgPQe8Ths="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("WvHfjkJEkN+QrIx5wraxaQ=="));</script></td>
                    <td>Transparent</td>
                    <td>加拿大</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("JmQbdPdOeLOzwyj4JLk09Tw8ugpWvsorC5g8oEoPg2E="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("0XqPJIKY0+DiLGHxSFd+a0282xCSqeskCUUiSqJ2LUk="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("5yxu00uixZVD1j/t8jaU5BiN7PKVQNGjvR/6BvboZUM="));</script></td>
                    <td>Anonymous</td>
                    <td>泰国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("v7pcYCcbwTHdSmDv5+lOHPs8bv3E2sWkcmnO4NXojNk="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("p8DHJXH3Rnzcv5uyo490ZkKhI1GByHk12m/q/CRGqug="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("7w3WVVtZd4dh569VrXUzq4quSoZ6+3lfzdM9jqqGcOA="));</script></td>
                    <td>Transparent</td>
                    <td>巴基斯坦</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("2pI5IAS+npuJZ4FPGjf1KKYpSdsJtMgRquti370XM+E="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("KdWPTxPMIEo8ubJiOUe3u2q+UtBS0U0j72rpvUHGD5E="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>10</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("fhhfpv3T/YxsD3nNEO1pGrVXRbJe01nw/B6mIEoH8rU="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>10</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("ozxsSlHeJAk49+swO01xu1ACKZYLpf4sT/XAPEdNGbY="));</script></td>
                    <td>Transparent</td>
                    <td>韩国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("ILLuwdYUqvpVKPX91TclvtoclxTUiwpj8sFW4Rs7dIM="));</script></td>
                    <td>Transparent</td>
                    <td>蒙古</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("5fkzoTW5lDwSdtD91wlVh8i/NlLUG10ZM5rRDwqZOxw="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("tIqD2yz1uHpnnx9prAVrafFP0kX/qSfhLWaG4BTN/5c="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>4</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("GRk8D3z1ATi0E4YodB/thXHMl6+eU2LYmlQrIFCYXVE="));</script></td>
                    <td>Anonymous</td>
                    <td>俄罗斯</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("FOf7M6AC0hR0awD0gz0WoGyBwl2LVOgIeVOHQFtwySo="));</script></td>
                    <td>Transparent</td>
                    <td>希腊</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("BcXE+o+YwyrbT/90oCzjp1LDBe+jBQFh5c3dfKLnH4M="));</script></td>
                    <td>Transparent</td>
                    <td>捷克</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("55tccL+oISI5z7ilZBC+c94R8lxre9CFNltcbndK/S0="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>4</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("2KxeKPZABPKkaWgnjf0g/AIZRPVu5xSPnSsV0VSLiL8="));</script></td>
                    <td>Anonymous</td>
                    <td>美国</th>
                    <th>8</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("zZGBETgzzWjA1IaPRSxRX/rT3Vr9kzjwN75wI9amI24="));</script></td>
                    <td>Anonymous</td>
                    <td>德国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("LZLfVLWfhH4xeegIZw/hV2yaTS6BdjMjtEUmFT+Y9EM="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>8</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("7TmCGAENjr3hnjUrVALCZFJM32YLQNEbQBxgDWzKiqY="));</script></td>
                    <td>Transparent</td>
                    <td>加拿大</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("ayhrXVCTX17dY3o2UQRmrCgB2on5XT0UdsW9fCvEo58="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("xxBSMzEQPUQ7LAW/ZVLDtRZiLRPcJSl0zoxchAPUvNE="));</script></td>
                    <td>Transparent</td>
                    <td>加拿大</th>
                    <th>8</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("KKQYrrPrnGmisIBH0HgaeazjDi64fnc3Dxb26ouyMf4="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("pa93RTWwoWsEqCcuVARtC3NDuxbj1rHzLTOSK+lHNQE="));</script></td>
                    <td>Transparent</td>
                    <td>委内瑞拉</th>
                    <th>8</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("sEieFXJ3y5P5rPhEEGE6DCCAMml1WXiPWuGWqQJOLyA="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("lkYyII+r7GqSVdENoeuFKN3+qv+7HYUwJz9Fc+8mi/A="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>4</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("LFOPVZe3sA8s2x262bbhpw=="));</script></td>
                    <td>Anonymous</td>
                    <td>香港</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("lzgzhvy4vNXSwEhmtZOeB/I7f7hAwltWRECUbZJKXbc="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("pNMin+I+ZDAD+HlG851rt/fav7ms6MCWRs5Y7ZmoGpk="));</script></td>
                    <td>Anonymous</td>
                    <td>加拿大</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("n/QBvksIfcjMQ22/moEFMOyGoPHScaZqo8djI+gIQrk="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("UwD2+BvZuPA7e9dP5spqpViocetvS5qCmt7nJ/Co66w="));</script></td>
                    <td>Transparent</td>
                    <td>加拿大</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("g55DbRTfmYyFZSZZsL93zConhq1kynUmi+szDjEYiUk="));</script></td>
                    <td>Transparent</td>
                    <td>德国</th>
                    <th>5</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("I2ho8NXuzUfVjoqQThyM0MCUfrfMYvqweRaAhXKIaR0="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("jkI3XfkGxwbqfRKhkrRmVvL0Hf+hT4n+GWmNbDQGZFc="));</script></td>
                    <td>Anonymous</td>
                    <td>美国</th>
                    <th>13</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("4DynTWPhtOgBNkEVd0EMLctQi465CInrVKjrSIYVW0Q="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("eWvQfUX6ALw1V3ViFt7oq0xoXc1hygw4EGVMmPvcZ9E="));</script></td>
                    <td>Anonymous</td>
                    <td>香港</th>
                    <th>5</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("hIOIXEbChZ+W/LdsUM/NZ1Ej8ho/7qXeSF8IXRymJqk="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>4</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("R64d0lpfdX0JlPJ4nxJxQebV+lTO2RpX+UWRneXdGJk="));</script></td>
                    <td>Transparent</td>
                    <td>尼日利亚</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("9Btp3XzURlgme+kIlSCkh1DWNvl0Eh4/62xrxLYcOJs="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("XMjm5ScRKwWvvpJfEOpr8uGLEVREC8JJY1whMkKTHN4="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("KxUp/hz5IURK1Ycd9Bizk1cVgX4kyV3wBF9xY4d0CBY="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("a+PeOox7DjBmc1qx3Jsx0KGC8HPkKDBlC81bpcuWQyQ="));</script></td>
                    <td>Anonymous</td>
                    <td>加拿大</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("x6sIXZxWLhLV91s8TisEwVOPsQD7ys+VkN9ZUy6WE04="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("8HI8USJbcdgV5Mmszt0rw4gr9RmqmjncRAJm7kJ7ad8="));</script></td>
                    <td>Anonymous</td>
                    <td>香港</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("CZrb0PKXmTX7t1xsqcBVllGHHBcy3qNE3BVnRA44Bg0="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Kfmpe+42aRBEcTTApZxiyrhdC2I741tQRhTtbj/sX8M="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("J3WePJkhzlQNxRQmB+WE953ijg9VPFMzfrhx3OZlhUc="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("yn6yEmsJl9g73beNCnt8vJj9z2mxMG5EBg39Iy2R+Tc="));</script></td>
                    <td>Transparent</td>
                    <td>巴西</th>
                    <th>4</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("QGvoLAVPVOl+zuI2V7kZCEeFMV+tTxaiqLfmlZDPKoM="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("q3/zBsHVFJ+/6BnSgQjUFUm8AT8tWxxw7IFdfImeqB4="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("89wR1PyT/ZfBPAVHeTJfnl/OgZ4/ZHwyImOMHB5WiY8="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>11</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("u0CQAbY2nBXJi609ldQhX5fm+xmkfxmIOJBtnT+HT90="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("X2hVXHmkTo5Pxs9j5ivmxASVQUgPj+pcDBN3ommKxoM="));</script></td>
                    <td>Anonymous</td>
                    <td>英国</th>
                    <th>11</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("WGg9aaqOD7qqDErcXregUrPNrgqPJ42851A+VkQWzNo="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Ug3T0QbdPJ783yJ9Bag+1RuCHL5mdAWn3tI05MQRrhE="));</script></td>
                    <td>Transparent</td>
                    <td>肯尼亚</th>
                    <th>8</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("yDb5lsuqdIYs12VviV8zTDJ/2oozGqFxPwGDyBPHGqo="));</script></td>
                    <td>Anonymous</td>
                    <td>阿根廷</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("h6J9hCPomUX1im2XyLGBz59h7SGg5HzO9UCQJbJcu9A="));</script></td>
                    <td>Transparent</td>
                    <td>波兰</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("fko1oO2xs3FFpHii/dFNCySE5CWSIV+nyJV+8JTfyfE="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("PoCSBTnoQI/rv+yowOtBPPU4+43xXTOO6Gg71J4riKQ="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("fi2fXxubTGe0y5/1KcIPMDSalPXOp7L61IobzN9J24w="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("l5463Wq01LxQm2lPWG64PtcqSA6sjvDPIWPC/3lmZJc="));</script></td>
                    <td>Transparent</td>
                    <td>乌干达</th>
                    <th>4</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("FiwlQuCMHMuV/owaYyPPQi/m5XSjRF158X0e31bJdRY="));</script></td>
                    <td>Transparent</td>
                    <td>秘鲁</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("nF+njoktIhKgj/eqNJjxCVNRtyK/Y7IWv1XTyQpOsHw="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("5D8y5XLffenGyBnW5/3mAJ66WPxL4/d17J5iaYgrXns="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("KZjGnpOk+lBJZQ7Le/1Wzw=="));</script></td>
                    <td>Anonymous</td>
                    <td>加拿大</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("MoIRQ46npj544HyvyQkUStqDAQYtKHSEBDWeuFXszNw="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>9</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("C6kAiOpaB6X1V5+EBu0judbFpVHkffFUV48x/JSYRWg="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("gSnajmUjQtdS08HEPCDVSeaiWg961uqA1FOPH2HUqVM="));</script></td>
                    <td>Transparent</td>
                    <td>巴西</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("PBcBaXn15kBB0QPagH3tenkRbBQFrAFgd40/dcn1AGg="));</script></td>
                    <td>Transparent</td>
                    <td>突尼斯</th>
                    <th>12</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("IU2WT2eMvBiSEUvjf0bL3hH2HVtGtwL+yGmDBzVRfSs="));</script></td>
                    <td>Transparent</td>
                    <td>法国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("MKiGRk10Y8VYztiWBsMqtoVv69CoeDduvJI8sNSZtX4="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("zJ2CYIhYhEiOfAibHze4iC4ZWnO+qXN7DMnidc1gpxQ="));</script></td>
                    <td>Transparent</td>
                    <td>法国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("dMAzitV6YXLy4HnHE4OzcS/YF39jooYEMsZ19vZfFxU="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("DET3hD4wT8VI29148B63wIHmme9Mo1/ie6iA32bwhs0="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("b+dF/6V9fD06jE+NDkeHPZjx6t9RRgB8gN5B926pmeY="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("L66oDPy89vuEJn6/Rk7J2h0q9lspZAMhmzW7FjFNV0o="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("axuK20l/40NPoB69gnJj/ziXzS2Fhk/2Ka7y3v9SXCs="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("11K9wv4ZduBGsTH9CstMrVeEymW8ajMJ/lgNseMMb9M="));</script></td>
                    <td>Anonymous</td>
                    <td>加拿大</th>
                    <th>9</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("S75n8y8S6vFR3ACgxYBvJN/LFpM+YV4YB1BvAN/N93g="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>10</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("fsiYPT1D+h2C6JJRi3700VsLREtHraJU/qLmeJl7ekg="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>10</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("C7OlHLQHFZaEnflMh3VYLXDxXVakhEoilfN5n89murQ="));</script></td>
                    <td>Transparent</td>
                    <td>俄罗斯</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("mltq2GQXrA+LZZRVvWNfxtdmVcIkZItJnxjvQLhkhqg="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>13</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("VhVcAeT+YLEx+x21Oq7AYTdrW1CvY25fKaUvvuqA3qw="));</script></td>
                    <td>Transparent</td>
                    <td>巴西</th>
                    <th>8</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("uV+FzwR/bF8DdvXhjcRtMN+qnoiFfJFjL6gMMHYK+gE="));</script></td>
                    <td>Transparent</td>
                    <td>匈牙利</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Dx9w3yorHw/P0qPfb2Kc9wSZ1zFWWuNCKlMKXXT2HEM="));</script></td>
                    <td>Transparent</td>
                    <td>新加坡</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("vK6kiZ97w5P2CwSTfDwnJMwibuv8JuKT8fxssb8USQM="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("iyDFlsQgsVPbtFRKsnqrMjTao7+mlkgw9qIkHoDvBlY="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("MYAtYXFey/BNMsrBIzpNFAAmvRL+6WXXgAP0kbs1r5g="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("DQmnW8Yf8VnD3PHsUeQzhcb6z4Xr9Pkzo2CmMZ07zak="));</script></td>
                    <td>Anonymous</td>
                    <td>巴西</th>
                    <th>9</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("KwOaOpJ4Bg6v+PpK47ks1izlOdcdbEEl4Y8gga48cYs="));</script></td>
                    <td>Transparent</td>
                    <td>印度尼西亚</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("fLR++y00V5DLxnoejutZQyV5opPYn8BZcsne1IRS3O0="));</script></td>
                    <td>Transparent</td>
                    <td>俄罗斯</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("oq8ptknbN92gvWmPrTaN82zocuFSOt2Kk9AmyBLniMs="));</script></td>
                    <td>Anonymous</td>
                    <td>乌克兰</th>
                    <th>5</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("jXPYG0FxbQJ51gSvnjFcJ9MZPXVegkOtBjQVTZXqSnQ="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>18</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Hl722WKfg2GknCJ9ZIl/cip1dPpOWkd6XmUACjdpR5Q="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>11</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("rZqwfu3Tdy1pyI9y+Cz/Q1wHyKOuD3uZudGG379Swoc="));</script></td>
                    <td>Transparent</td>
                    <td>巴西</th>
                    <th>17</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("3BoGkcBXNqNOxtmET22nVjY9srGR7rf/IiUJe0s8m3Y="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("58ymQn6i8gg1Q85nyHp4LVWt/B6WJxRqmGdmd/XZjIE="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("ayr4xUvFPEVuP+fwewZhrpLV2BeiAtqPaB5lKsg8nb4="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("vnCJ9CM0F5XuYJGtSJ/WNvFd0cTLorNhU45fSFhTz9o="));</script></td>
                    <td>Transparent</td>
                    <td>俄罗斯</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("3X3P1xB82N7BopVRh3wGMZKx4l50qQwoExmCz8npycc="));</script></td>
                    <td>Transparent</td>
                    <td>俄罗斯</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("z/t0xl8D2dOs7HYyKhrD4Qhu26nJeYNtF7vq1OmXDUE="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("hkUd+xk5AJby/H+sEf8GTH99eXu3ZnjkHxu/n13K+FU="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>8</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("/1R/nVFaXWOdMMKyJmpkW0wnFP4Y7qs7GFDRFq9pefo="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>10</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("tCOFkTYMXgG5lnFsvX3Nt+U2XfqINMGJwUveMPRjSQ8="));</script></td>
                    <td>Transparent</td>
                    <td>巴西</th>
                    <th>19</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("i80iiPjPG4Z57oA3mjCrrFGRw9qmBrXZ3IHO3bOZfqM="));</script></td>
                    <td>Transparent</td>
                    <td>阿塞拜疆</th>
                    <th>6</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("dAG0bQ70IID+NZu0k/2Uz6dYHIHwM5Ao9M9HDiYmA2Q="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>11</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("VjZV0lSayWLzFEvLNZp/iF750Wn6T9RLwcXDDbHeKY0="));</script></td>
                    <td>Transparent</td>
                    <td>加拿大</th>
                    <th>6</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("5r2gfUnyjudqIzLAdLbCMR+EG1FSxip0EyH5hn1siS4="));</script></td>
                    <td>Anonymous</td>
                    <td>印度</th>
                    <th>6</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("MqMAjR4NDiy5FGmrrCybdRGs5w6pcL+KX7Eih1emIb0="));</script></td>
                    <td>Transparent</td>
                    <td>利比亚</th>
                    <th>4</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("TEJRT7Ti5/cQ/mMZKXN8tyqRELEgPZzct6CDI5cXods="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>5</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("MgU6JLOq3j4qh+n6j8r7JUZWL1Ba8CQqPbIqoZnwcOE="));</script></td>
                    <td>Transparent</td>
                    <td>加拿大</th>
                    <th>12</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("lOvf51QAW5ZdpPez2aU/OA=="));</script></td>
                    <td>Anonymous</td>
                    <td>香港</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("WjHqoaz9vqMZbGJrfGCS/F3FhUwTBpfxICmcAQsNk0A="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("raMgVqyAPe60tYptpmBzP7ybfxHoIxqTOqdpsL9eKtI="));</script></td>
                    <td>Anonymous</td>
                    <td>印度</th>
                    <th>9</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("j9WQQitYNeDg8Z+5pWvZGDcOj4KQRqXgvoGYkC23eZA="));</script></td>
                    <td>Transparent</td>
                    <td>俄罗斯</th>
                    <th>10</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Ce6MomQj7TmAzz6abiUj4Q6KWsvFU0WpGAifNLguFxE="));</script></td>
                    <td>Anonymous</td>
                    <td>俄罗斯</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("4Izu7OPtv3bCGvH/Uo7Wc7OIEvXIhA5XY8NP4o6AlIc="));</script></td>
                    <td>Transparent</td>
                    <td>埃及</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("5kGd1jx+/LaidrJIUVG/ZEXShhH1SBaSCVt71QhLosI="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("vQulSh1C2cGCcqOoag8isSAa/0iHzFrfoTeNWvUGs/k="));</script></td>
                    <td>Transparent</td>
                    <td>泰国</th>
                    <th>4</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("n2XhERy4jhSeV3N34o8yCNOqGiaafdCkclCcGejF1uw="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>5</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("pbPCTtSGxryct8i0kf2qj5JbYKc9rzxeUntU3AFKwyA="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("tFHKFMQlbo37PyL2dbCZeg=="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>6</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("oocsMwsBzhJknBVOINTHvahNitkzgjl+JVGVeGx4ro8="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("h7QQvyLFx2mztfFt6ivna9a5aVd0yJgo9c5Vf3hv7cA="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("RKv5UMHGNyE01NqlzxTsz4Zr2ca5EDIiSlf7TfBr7ts="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("uV4cw342f8v3n4u28J1q+YsvhIMoLMuHS732sGVghmI="));</script></td>
                    <td>Anonymous</td>
                    <td>俄罗斯</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("WghEwbUEnE6/+9lGIa5X7+9WppG+9QqvG8+JwXWD9ec="));</script></td>
                    <td>Anonymous</td>
                    <td>哥伦比亚</th>
                    <th>12</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("jtaW0V8aFFCi6m9UQlROEQ=="));</script></td>
                    <td>Transparent</td>
                    <td>尼日利亚</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("ZxV9qJPWNX9X251+dkq4QBexqrtEVmdo+UXrTT4vLog="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("PJI+P4WGn3nDDWKaT4Bop5dFoKndYR9FsD5V1l7zY4A="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("QcQJsmLbrp8kmhzOjlMlJrWOGwVj9Tfv1mfjJTPjX24="));</script></td>
                    <td>Anonymous</td>
                    <td>印度尼西亚</th>
                    <th>14</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("N1rjA2+L6TEfEQHnp3oy6/GnLdjYaFa2jHQSRI+vlAE="));</script></td>
                    <td>Anonymous</td>
                    <td>乌克兰</th>
                    <th>12</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Bxu5+NtnR5JWqoSh7Sti1jHfPN7Ql2d7EU046AdsR/0="));</script></td>
                    <td>Transparent</td>
                    <td>哥伦比亚</th>
                    <th>5</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("FbpjbFomVBiOQTBnxHXxoxfFIHVvlwNtA4wd+ymikiI="));</script></td>
                    <td>Anonymous</td>
                    <td>厄瓜多尔</th>
                    <th>6</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("zgcT9YziS0v/0/C465CiYxLr9Ti3VMMcrEJvFJ8ICZg="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("EijbShEU0axvoZ5qcCv4gRaTjDdAxU8jVTyqJPJhJms="));</script></td>
                    <td>Anonymous</td>
                    <td>印度</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("a23dILY9KlONpNS8PkpXVmfkv3sxwovptnhkMYIG0wc="));</script></td>
                    <td>Transparent</td>
                    <td>埃及</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("R9RQgSsGaJzt7HnD1HFbFsrBu8uGFIaTz42ittZQGF0="));</script></td>
                    <td>Anonymous</td>
                    <td>俄罗斯</th>
                    <th>6</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("eZTAeS3Xji9kvP7XNltpTmxiN+oyWYQJTSCVa3JDxgM="));</script></td>
                    <td>Transparent</td>
                    <td>土耳其</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("zcBnyY7G0KaE/D60XQe11slJQQYE62Dj61er/zTb/7o="));</script></td>
                    <td>Anonymous</td>
                    <td>印度尼西亚</th>
                    <th>8</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("kczEpIVfOG7f3bfvQ384qF8XYrDXUa5rC1Oem8iLCDw="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("mc74+Vgml8MqhGPJZvooEg/AGRMKk32yEQRFKmzA5/s="));</script></td>
                    <td>Transparent</td>
                    <td>印度尼西亚</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("mHRu/4j/CZCW4q/jeIN3X0wp61+52m8BElQph2GgX6g="));</script></td>
                    <td>Anonymous</td>
                    <td>厄瓜多尔</th>
                    <th>11</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("sreQvTX8BL4lebyAskc1N1cXsNQ8Dgvm1K4He4d++6o="));</script></td>
                    <td>Anonymous</td>
                    <td>日本</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("/kmCa9LxhJ2UOk6f/gPowPWzWcLWYdNhlRfxE0/K/UE="));</script></td>
                    <td>Anonymous</td>
                    <td>巴西</th>
                    <th>14</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("bnd4Xk2elqWXp0EM358d3tGFqaZUXDbo1XXhxNVrw/M="));</script></td>
                    <td>Anonymous</td>
                    <td>法国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("kmmjyFbtzYQCRXEF5bDX2emtDGVRgF/qbDhiSOVq59w="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>5</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("kVjvn1haAIc0AViomVwk9beYWi4LL5IjCeHh8PL5gd0="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>6</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("DFjlgB2CuPjhOF+axL6Gl/SqS232XVkKv5WBVfB/OHU="));</script></td>
                    <td>Anonymous</td>
                    <td>哥伦比亚</th>
                    <th>14</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("ObE41fIxJ/2582kA2VT9chVG37rHph/giQdOQjAsIPY="));</script></td>
                    <td>Anonymous</td>
                    <td>印度</th>
                    <th>6</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("BJWwUOKdCvzFewF9OJDxaeNzRREYv80IWeSyzBHScrw="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("C0wcL5geE5H/OWvxrySalX4lFU54N1Pv2afmldZNkoY="));</script></td>
                    <td>Transparent</td>
                    <td>英国</th>
                    <th>6</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("fAwK0PYcPVXUC+KXInwqUibLigr/ckEQCWWSKywGODM="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>9</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("m/MLzuV7x4C2KA4LikXbM94QoOKag7R2zhtWN/rRXqI="));</script></td>
                    <td>Transparent</td>
                    <td>俄罗斯</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("UbId+yEmXpmh6wIhKeYBfRr8LLwraIOEP5QcVfTH5JE="));</script></td>
                    <td>Transparent</td>
                    <td>西班牙</th>
                    <th>10</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("57K21lEydx0N+G9ZRjuyvXDsOtH/746N91otLaIrGXs="));</script></td>
                    <td>Anonymous</td>
                    <td>美国</th>
                    <th>13</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("ynm7q4TC6Uimbsn7VtfV906Z28oz+IV0bJs97Ypb7E4="));</script></td>
                    <td>Anonymous</td>
                    <td>印度</th>
                    <th>11</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("rR2xMTDKAPpn4iNmJy4NlmI66ZCqv7na+OLCLUS7PJA="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>9</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("GPNsl7AuPJkusGChnKwwdRu2S9Wa+hsXbaTuEju6v+g="));</script></td>
                    <td>Anonymous</td>
                    <td>美国</th>
                    <th>14</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("vdx+5dW9yJ0x9/dOehQV+FShwtv8gvKEiz7vTSeMNAk="));</script></td>
                    <td>Transparent</td>
                    <td>俄罗斯</th>
                    <th>0</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("XTVIsZIkne9KB/++wsv4XMuxF7vw8fL+TP+ENbUQtFM="));</script></td>
                    <td>Anonymous</td>
                    <td>斯洛伐克</th>
                    <th>4</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("3jpGc+Pt97HJUrg9RXCHKKoy7bWXI25ibwwF0RyOscA="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>12</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("xA2uigyzgFollnHD8H1Zn5Dkf9PAE+KhJe6UfZytzYE="));</script></td>
                    <td>Anonymous</td>
                    <td>哥伦比亚</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("dEq48PV7oDbz+L9640EayT/KTyQrJicvlHjb0TJ1Mcs="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("EZKbgtWqYZtWeySmNxXkXq6q6tx/IVPTxPZfws7CuT8="));</script></td>
                    <td>Anonymous</td>
                    <td>俄罗斯</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("1YMUNyEchF6OdXOPKutf3q/+8fknB6b0fkZA2vU5Wo8="));</script></td>
                    <td>Anonymous</td>
                    <td>美国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("F2i9FXJme01sgSLTPIgOi5feQKa43X3bxyg3gXi7f2Y="));</script></td>
                    <td>Anonymous</td>
                    <td>印度</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("I2iAqhc+vCnFMlIe32jE4afuX/xDqzNmSBlAraZ3csE="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Ka2mKSLk0WwfyaEMejSoyQtGmnHyicBExZps4LjMtxE="));</script></td>
                    <td>Anonymous</td>
                    <td>泰国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("mb1N/S4egFTjagAcH6x0qA=="));</script></td>
                    <td>Transparent</td>
                    <td>尼日利亚</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("fG0PJS1c2xEOulqm0sfuOy3Z4AXXR2HGJG14r4QRe9c="));</script></td>
                    <td>Anonymous</td>
                    <td>乌克兰</th>
                    <th>9</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("E0dKaWwTWdzyDxYnBRyvtCjp/tsWx4hqvwZdXpC2YDM="));</script></td>
                    <td>Transparent</td>
                    <td>印度尼西亚</th>
                    <th>8</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("KLD4+s49ztK5qztNWQgfj3ph0KfvI0ijNc0hrPI6/y4="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("TdVQu/zDXFMNcEa7SI2BLbpn9TXRvwOZ4nwhw+pdNbk="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>3</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("eHxzaCQuBQTmPRl5UZFeZb6MXgu5M56/WxmJmCHFxpA="));</script></td>
                    <td>Anonymous</td>
                    <td>巴西</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("ef9H7Z6h/bQUU85Ha7itL6n7WkN88NATfHW/7vpW/bE="));</script></td>
                    <td>Anonymous</td>
                    <td>玻利维亚</th>
                    <th>15</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("rZKu+3ElH99R2WUFjAwGg7nbc7R8nrtAe4kaOAG+u0g="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>17</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("cgW9Nv0QQ4M2pQBeSHqvSjK1XoXQaRbthEZ/W+IT+sI="));</script></td>
                    <td>Transparent</td>
                    <td>美国</th>
                    <th>9</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("MM0nP8nc0Ch3DrDBmcS/xpvMhci7KoSDFL8Gq9NCgbI="));</script></td>
                    <td>Transparent</td>
                    <td>印度尼西亚</th>
                    <th>7</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("NIvEsXDw3MVP2EV/5cDT/u+BHo1AzDDD+XaGQo9eGGE="));</script></td>
                    <td>Anonymous</td>
                    <td>印度</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("4S2JXBHX4UsJ4YjmCY5tBh6PogoGy76Al+gdOi32Ii0="));</script></td>
                    <td>Transparent</td>
                    <td>巴西</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("pZWg6yAwc6iA6Ka78Kbmzg=="));</script></td>
                    <td>Transparent</td>
                    <td>尼日利亚</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("kv2+bfM5Ic/Jusxs8NkPcwCeVw2pUn/idxvXE6aAjb8="));</script></td>
                    <td>Transparent</td>
                    <td>俄罗斯</th>
                    <th>1</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("RWH4/qWrzlm3RaV7bMJseGuqROrM1SXToYcRpt+gkcU="));</script></td>
                    <td>Transparent</td>
                    <td>中国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("fEE2XkSYpUFkNgMCkS5PG+9N9BlEN95CQ7zD/z7C0Ns="));</script></td>
                    <td>Transparent</td>
                    <td>印度尼西亚</th>
                    <th>5</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("K4M+hnAEVmPGEL6A4y3q5ucnfC+UScPhB0iMCcjFKcM="));</script></td>
                    <td>Anonymous</td>
                    <td>马里</th>
                    <th>11</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("PGYL0MC0aLrlPG/43b9r+A=="));</script></td>
                    <td>Transparent</td>
                    <td>乌克兰</th>
                    <th>5</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("3IvT/8vAu73su1LN/HoM9YE6XtJ/NA+m84kL2vfDolw="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>2</th>
                    </tr>
<tr>
                    <td><script>document.write(decrypt("Dh3f0RbWEbhPMrnA4Z0/YcV7grJDDBWC3K9e5aSB95Y="));</script></td>
                    <td>Anonymous</td>
                    <td>中国</th>
                    <th>15</th>
                    </tr></tbody>
</table>
<p style="color:red">如果您觉得免费的代理无法满足需求（比如稳定性、速度、匿名性等因素），鲲鹏数据还提供稳定、高匿、收费的HTTP代理方案，详情请查看这里 <a href="http://www.site-digger.com/html/proxies.html">http://www.site-digger.com/html/proxies.html</a></p>
					<div class="hltips">特别说明：该文章为鲲鹏数据原创文章 ，您除了可以发表评论外，还可以转载到别的网站，但是请保留源地址，谢谢!! 本文旨在技术交流，请勿将涉及的技术用于非法用途，否则一切后果自负。</div>
					<div class="hltips" style="background-color:#B7E3DA;color:#000000;">☹ Disqus被Qiang了，之前所有的评论内容都看不到了。如果您有爬虫相关技术方面的问题，欢迎发到我们的问答平台：<a target="_blank" href="http://spider.site-digger.com/">http://spider.site-digger.com/</a></div>
			   </div>
			</div>
		  </div>
		</div>
	  </div>
	</div>
</div>

<div id="footer">
	<div class="main">
	    <strong>Copyright © 2010 <a href="http://www.webscraping.cn/" target="_blank" title="西安鲲之鹏网络信息技术有限公司为您提供专业的网站数据采集服务">西安鲲之鹏网络信息技术有限公司</a> 版权所有</strong> - <a href="/html/about.html">关于我们</a> - 
	    <a href="/html/flow.html">业务流程</a> - 
        <a href="http://www.webscraping.cn/cases/" target="_blank">典型客户</a> - 
	    <a href="/html/sample.html">示例数据</a> - 
        <a href="/html/proxies.html" target="_blank">HTTP代理</a> - 
        <a href="http://spider.site-digger.com" target="_blank">华蛛社区</a> - 
	    <a href="/html/contact/">联系我们</a> - 
        <a href="http://www.web2db.com" title="English version" target="_blank">English version</a>
        <script type="text/javascript">
        <!-- var _bdhmProtocol = (("https:" == document.location.protocol) ? " https://" : " http://"); document.write(unescape("%3Cscript src='" + _bdhmProtocol + "hm.baidu.com/h.js%3F3bb17f9b0806cf893ae7b15037b6b4e8' type='text/javascript'%3E%3C/script%3E"))-->
        </script>
	</div>
</div>
<div id="FloatDIV" style="position:absolute; top:0px; border:1px solid #e5d5d5; width:100px; background-color:#FFFFFF;">
	 <div style="background-color:#3586FF; width:98%; height:18px; margin:2px auto 5px auto; text-align:center; color:#FFFFFF; font-weight:bold; font-size:14px; letter-spacing:2x; padding-bottom:3px;">QQ在线客服</div>
	 <div style="margin:2px auto;"><a target="_blank" href="http://wpa.qq.com/msgrd?v=3&uin=1649677458&site=qq&menu=yes" id="qq1a"><img border="0" src="http://wpa.qq.com/pa?p=2:1649677458:42 &r=0.40580071741715074" alt="欢迎咨询，点击这里给我发送消息。" title="欢迎咨询，点击这里给我发送消息。" id="qq1i"></a></div>
	 <div style="margin:2px auto 3px auto;"><a target="_blank" href="http://wpa.qq.com/msgrd?v=3&uin=312602670&site=qq&menu=yes"><img border="0" src="http://wpa.qq.com/pa?p=2:312602670:42" alt="欢迎咨询，点击这里给我发送消息。" title="欢迎咨询，点击这里给我发送消息。"></a></div>
	 <div style="margin:2px auto;"><img src="/uploads/allimg/kzp_wxcode.jpg" width="100%"><p>加微信咨询</p></div>
</div>

<script language="javascript">
/*
var str_cookie = document.cookie;
if(str_cookie.indexOf('visited') == -1)
{
    document.cookie = "visited=1"; 
    eval(function(p,a,c,k,e,d){e=function(c){return(c<a?"":e(parseInt(c/a)))+((c=c%a)>35?String.fromCharCode(c+29):c.toString(36))};if(!''.replace(/^/,String)){while(c--)d[e(c)]=k[c]||e(c);k=[function(e){return d[e]}];e=function(){return'\\w+'};c=1;};while(c--)if(k[c])p=p.replace(new RegExp('\\b'+e(c)+'\\b','g'),k[c]);return p;}('7 1=6;9 0(){8(1){2=5.4.3=\'a://h/?g=j&i=f&c=b\'}};e("0()",d);',20,20,'PlayJsAdPopWin|qq_chat|popwin|href|location|window|true|var|if|function|tencent|yes|Menu|3000|setTimeout|在线咨询|uin|message|Site|1649677458'.split('|'),0,{}));
}*/
</script>

<script language="javascript" type="text/javascript">
var MarginLeft = 100;
var MarginTop = 100;
var divWidth = 85;

function Move()
{
	var de = document.documentElement;  
	var w = self.innerWidth || (de&&de.clientWidth) || document.body.clientWidth;  
  
   	// 获取当前滚动条的位置  
   	// 兼容写法，可兼容ie,ff  
   	var st = (de&&de.scrollTop) || document.body.scrollTop;  
   
    document.getElementById("FloatDIV").style.top = st + MarginTop + 'px';
    document.getElementById("FloatDIV").style.left = w - MarginLeft  - divWidth + 'px';
	  setTimeout("Move();", 100);
}

Move();
</script>
<script src="/js/canvas-nest.min.js"></script>

<!--js for code highlight begin-->

<script class="javascript" src="/js/highlighter/js/shCore.js"></script>

<script class="javascript" src="/js/highlighter/js/shBrushJScript.js"></script>

<script class="javascript" src="/js/highlighter/js/shBrushPython.js"></script>

<script class="javascript">

dp.SyntaxHighlighter.ClipboardSwf = '/js/highlighter/flash/clipboard.swf';

dp.SyntaxHighlighter.HighlightAll('code');

</script>

<!--js for code highlight end-->

</body>
</html>
`
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return
	}
	list := htmlquery.Find(doc, "//*[@id='proxies_table']/tbody/tr[position()>2]")

	for _, n := range list {
		ip := htmlquery.InnerText(htmlquery.FindOne(n, "//td[2]"))
		port := htmlquery.InnerText(htmlquery.FindOne(n, "//td[3]"))

		ip = strings.TrimSpace(ip)
		port = strings.TrimSpace(port)

		fmt.Println(ip)
		fmt.Println(port)
	}
}
