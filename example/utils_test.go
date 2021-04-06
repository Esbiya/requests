package example

import (
	"fmt"
	"testing"

	"github.com/Esbiya/requests"
)

func TestRandomUserAgent(t *testing.T) {
	//ua := requests.RandomUserAgent(requests.Chrome)
	//log.Println(ua)
	//ua1 := requests.RandomUserAgent(requests.Safari)
	//log.Println(ua1)
	//ua2 := requests.RandomUserAgent(requests.IE)
	//log.Println(ua2)
	//ua3 := requests.RandomUserAgent(requests.Opera)
	//log.Println(ua3)
	//ua4 := requests.RandomUserAgent(nil)
	//log.Println(ua4)

	// a := requests.Payload{
	// 	"1": "4",
	// 	"hhh-==": "ggg==",
	// }
	// log.Println(a.Stringify())

	r := requests.Response{
		Text: `<!DOCTYPE html>
		<!-- 接口方式支付,中间页面 -->
		<html>
		<head>
		
		<title>支付收银台</title> 
		<meta http-equiv="Content-type" content="text/html; charset=GBK"> 
		<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
		<meta name="apple-mobile-web-app-capable" content="yes" />
		<meta name="format-detection" content="telephone=no"/> 
		<meta name="referrer" content="origin"/>
		<link type="text/css"  rel="stylesheet" href="/upay/wps/css/global.css" /> 
		<link type="text/css"  rel="stylesheet" href="/upay/wps/css/style.css" /> 
		<META HTTP-EQUIV="Pragma" CONTENT="no-cache">
		
		<link type="text/css" rel="stylesheet" href="/upay/wps/css/interface_wps.css">
		<script type="text/javascript" src="/upay/wps/js/jquery-1.4.2.min.js"></script>
		<script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.3.2.js"></script>
		<style>
				.ordinfo{width:94%;border-radius:10px;margin:10px auto;height:225px;box-shadow: 0 2px 6px 2px #ccc;;}
						.ordnum{height:45px;line-height:45px;text-indent:15px;background:#ff7618;color:#fff;font-size:16px;border-top-left-radius: 10px;border-top-right-radius: 10px;}
						.ordinfo ul li{margin-top:20px;width:90%;margin-left:5%;list-style:none;font-size:16px;}
						.ordinfo span{float:right;}
						.ordinfo span em{color:#f98808;font-style:normal;}
		</style>
		<script type="text/javascript">
		$("title").html("支付跳转中...");
				var CAP_CNL = "WECHAT";
				var PAG_NTF_URL = "https://wap.js.10086.cn/WSCZYLRESULT.thtml?ordno=F0BCB6E21DF03624F5F16E493F8A8E0A73A8756A170A726E3EBEB419B2D737BD&loginMobileAES=73A8756A170A726E3EBEB419B2D737BD&mobileAES=489EFF1F988A33E4D0B835273DF8027F&actualmoneyZFAES=06A0889555D932ED059C589C50C97E24&bankTypeAES=BB8C09BA0A84FE44EDFEE39687999590&ticketFlagAES=FF3B45660B219DD30C9459003540C9C4&bossDateAES=3BE13625B44A7A27A954C57F7D841066&isBindYxaAES=46411EE8A66F18B36B8A258402BE9749&isSHHYFlagAES=46411EE8A66F18B36B8A258402BE9749";
				var PAG_NTF_URL2 = "https://wap.js.10086.cn/WSCZYL.thtml";
				var tjtime = 70;
				var tjtimer;
				var ORD_NO = "202104064729155538";
				var WE_RETURN_URL = "http://www.js.10086.cn/upay/wps/service/tpfWePayCallBackDq.xhtml"+"?ORD_NO="+ORD_NO;
				var WECHARTGZH_URL = "http://upay.12580.com:80/upay/wps/service/doTfpWapPayment.xhtml";
				var PAYSIGN = "";//签名
				var APPID = "";//公众号id
				var TIMESTAMP = "";//时间戳
				var NONCESTR = "";//随机字符串
				var PREPAYID = "";//订单详情扩展字符串
				var walletUrl="http://upay.12580.com:80/upay/wps/service/WalletToSmsPage.xhtml"+"?ORD_NO="+ORD_NO;
				$(function(){
						$("#we_return_url").val(WE_RETURN_URL);
						var buttonType = $("#buttonType").val();
						if(buttonType=="0"){
								showModel();
								return;
						}
		
						tjtimer = setTimeout(function(){
								var flag = navigator.userAgent.toLowerCase().match(/MicroMessenger/i);//=="micromessenger"
								var reserved = $("#reserved1").val();
								if(tjtimer){
										clearTimeout(tjtimer);
								}
								if(CAP_CNL=="ALIPAY" || CAP_CNL=="WECHAT"|| CAP_CNL=="CMPAY" || CAP_CNL=="EPPAY" || CAP_CNL=="QUICKPAY" || CAP_CNL=="SPDBTPF" || CAP_CNL=="CMBPAY" || CAP_CNL=="WALLET"){
		
										if(CAP_CNL=="WECHAT"){
		
												//页面在微信中打开，公众号中所应用
												if(flag=="micromessenger"){//测试用！=，正常用==
														$("#WPS_WECHART_TYPE").val("JSAPI");
														$("#RESERVED").val(reserved);
														$(".ordpayway").html("微信支付");
														//判断微信小程序方法
														wx.miniProgram.getEnv(function(res) {
																if(res.miniprogram){
																		$(".ordpayway").html("微信支付");
																		var xappid=getQueryString("appid");
		
																		if(xappid){
																				$("#appid").val(xappid);
																		}
																		$("#isXcxzf").val("1");
																		$("#isGzhzf").val("0");
		
																		var ajaxData=$("#payForm").serialize();
		
																		var path='/pages/zf/main?'+ajaxData;
																		wx.miniProgram.navigateTo({
																		url: path
																	})
																		return;
																}else{
																		$(".ordpayway").html("微信支付");
																		$("#isGzhzf").val("1");
																		$("#isXcxzf").val("0");
																		doGzhAjax();
																		$("#buttonType").val("0");
																		showModel();
																		return;
																}
		
		
														})
		
														return;
												//页面在浏览器中打开
												}else{
		
												}
												$(".ordpayway").html("微信支付");
										}else if(CAP_CNL=="ALIPAY"){
												$(".ordpayway").html("支付宝支付");
										}else if(CAP_CNL=="WALLET"){
												//钱包支付
												$(".ordpayway").html("钱包支付");
												window.location.href=walletUrl;
												return;
										}else if(CAP_CNL=="CMPAY"){
												$(".ordpayway").html("和包支付");
										}else if(CAP_CNL=="EPPAY"){
												$(".ordpayway").html("苏宁支付");
										}else if(CAP_CNL=="SPDBTPF"){
												$(".ordpayway").html("小浦支付");
										}else if(CAP_CNL=="CMBPAY"){
												$(".ordpayway").html("一网通银行卡支付");
										}else if(CAP_CNL=="WEICHAT"){
												$(".ordpayway").html("微信支付");
										}
		
										$("#buttonType").val("0");
										showModel();
		
										if(flag=="micromessenger")
										{
										   wx.miniProgram.getEnv(function(res){
														if(res.miniprogram){
														}else{
		
																$("#payForm").submit();
														}
											})
										}else
										{
										   $("#payForm").submit();
										}
								}else{
										$("#loadingDiv").hide();
										$("body").css("background-color","#fff");
										$("#nosupport").show();
								}
						},tjtime);
				});
		
				function closePage(){
					window.history.go(-1);
				}
		
				function showModel(){
						$(".modelBody").show();
				}
		
				function finishPay(){
						$("#clickflg").val("1");
						$("#interSucForm").submit();
				}
		
				function giveupPay(){
						if(PAG_NTF_URL2!=null && PAG_NTF_URL2!=""){
								window.location.href = PAG_NTF_URL2;
						}else{
								if(PAG_NTF_URL!=null && PAG_NTF_URL!=""){
										window.location.href = PAG_NTF_URL;
								}
						}
				}
				//获取地址栏参数
				function getQueryString(name){
					 var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
					 var r = window.location.search.substr(1).match(reg);
					 if(r!=null)return  unescape(r[2]); return null;
				}
		
				//ajax请求调用后台获取公众号支付需要的参数（公众号支付）
				function doGzhAjax(){
						if(tjtimer){
								clearTimeout(tjtimer);
						}
		
						var ORD_NO = $("#WPS_ORD_NO").val();
						var PAY_AMTS = $("#WPS_PAY_AMTS").val();
						var ORDER_DESC = $("#WPS_ORDER_DESC").val();
						var CASH_CORG = $("#WPS_CASH_CORG").val();
						var WECHART_TYPE = $("#WPS_WECHART_TYPE").val();
						var RESERVED = $("#RESERVED").val();
						var tokenVal = $("#token").val();
						var isgzhzf=$("#isGzhzf").val();
		
		
						var surl=window.location.href;
						var hurl=surl.split(":");
						if(hurl[0]=="https"){
								WECHARTGZH_URL=WECHARTGZH_URL.replace("http","https");
		
						}else{
								WECHARTGZH_URL=WECHARTGZH_URL.replace("https","http");
						}
						if(WECHARTGZH_URL.indexOf(":80")==-1){
		
						}else{
								WECHARTGZH_URL=WECHARTGZH_URL.replace(":80","");
		
						}
		
						$.ajax({
						type: "POST",
						url:WECHARTGZH_URL,
						contentType:'application/x-www-form-urlencoded; charset=UTF-8',
						data : {
								"ORD_NO":ORD_NO,
								"PAY_AMTS":PAY_AMTS,
								"ORDER_DESC":ORDER_DESC,
								"CASH_CORG":CASH_CORG,
								"WECHART_TYPE":WECHART_TYPE,
								"RESERVED":RESERVED,
								"TOKEN":tokenVal,
								"isGzhzf":isgzhzf
						},
								dataType : "json",
								success: function(msg) {
		
										if(msg.GWA.MSG_CD == "GWA00000"){
												PAYSIGN = msg.SIGN;//签名
												APPID = msg.APPID;//公众号id
												TIMESTAMP = msg.TIMESTAMP;//时间戳
												NONCESTR = msg.NONCESTR;//随机字符串
												PREPAYID = msg.PREPAYID;//订单详情扩展字符串
												//alert(PAYSIGN+" "+APPID+" "+TIMESTAMP+" "+NONCESTR+" "+PREPAYID+" "+PAYSIGN);
												doGzhWechartPay();
								}else{
										alert(msg.GWA.MSG_CD+"："+msg.GWA.MSG_INF);
							}
						},
								error: function(msg) {
										//alert("系统异常，请稍后重试！Connection error！！！"+msg.response);
										//alert(msg.response);
										alert("系统异常,请稍后再试！");
						}
						});
				}
		
				function doGzhWechartPay(){
						if (typeof WeixinJSBridge == "undefined"){
						   if( document.addEventListener ){
							   document.addEventListener('WeixinJSBridgeReady', onBridgeReady, false);
						   }else if (document.attachEvent){
							   document.attachEvent('WeixinJSBridgeReady', onBridgeReady); 
							   document.attachEvent('onWeixinJSBridgeReady', onBridgeReady);
						   }
						}else{
						   onBridgeReady();
						}
				}
		
				function onBridgeReady(){
				   $("#buttonType").val("0");
				   showModel();
		
				   WeixinJSBridge.invoke(
					   'getBrandWCPayRequest', {
						   "appId":APPID,     //公众号名称，由商户传入     
						   "timeStamp":TIMESTAMP,         //时间戳，自1970年以来的秒数     
						   "nonceStr":NONCESTR, //随机串     
						   "package":"prepay_id="+PREPAYID,     
						   "signType":"MD5",         //微信签名方式：     
						   "paySign":PAYSIGN //微信签名 
					   },
					   function(res){   
								  
						   if(res.err_msg == "get_brand_wcpay_request:ok" ){
		
							   }     // 使用以上方式判断前端返回,微信团队郑重提示：res.err_msg将在用户支付成功后返回    ok，但并不保证它绝对可靠。 
					   }
				   ); 
				}
		</script>
		</head>
		<body style="">
		<section style="visibility: hidden;">
				<div class="ordinfo">
						<div class="ordnum">订单号码：<span style="float:right;">1595000022173911</span></div>
						<ul>
								<li>充值号码：<span>15805125684</span></li>
								<!--<li>充值金额：<span>1297元</span></li>-->
								<li>支付方式：<span class="ordpayway"></span></li>
								<li>支付金额：<span><em>12.97</em> 元</span></li>
						</ul>
		
				 </div>
				<div id="loadingDiv" class="" style="text-align: center;margin-top:50px;display:none;">
						<div class="tzzDiv">
								<div style="height:20px;"><!--  --></div>
								<div class="interDiv">
										跳转中
								</div>
								<img alt="" src="/upay/wps/images/loading2.gif" width="40" height="20"/>
						</div>
				</div>
				<div id="nosupport" style="display: none;text-align: center;margin-top: 50px;">
						<div class="fail">
						<p>
							<img src="/upay/wps/images/fail_icon.png" alt="Failure"/>
						</p>
						<p class="pay_fail">暂不支持该支付方式！</p>
					 </div>
					 <div class="pd0_15 mtb50">
						<button id="" type="button" class="ui_btn ui_btn_s2" onclick="closePage()">关闭</button>
						 </div>
				</div>
				<form id="payForm" name="payForm" action="http://upay.12580.com:80/upay/wps/service/doTfpWapPayment.xhtml" method="post" ><!-- accept-charset="utf-8" onsubmit="document.charset='utf-8';" -->
						<input type="hidden" name="token" id="token" value="1617685224997" />
		
						<input type="hidden" id="WPS_ORD_NO" name="ORD_NO" value="202104064729155538" /> 
						<input type="hidden" id="WPS_PAY_AMTS" name="PAY_AMTS" value="12.97" />
						<input type="hidden" id="WPS_ORDER_DESC" name="ORDER_DESC" value="订单:1595000022173911" />
						<input type="hidden" id="WPS_CASH_CORG" name="CASH_CORG" value="WEICHAT" />
						<input type="hidden" id="WPS_WECHART_TYPE" name="WECHART_TYPE" value="" /><!-- JSAPI -->
						<input type="hidden" id="RESERVED" name="RESERVED" value="" /><!--  -->
		
						<input type="hidden" id="spcrip" name="SPCRIP" value="116.22.58.119">
						<input type="hidden" id="detailurl" name="DETAILURL" value="wap_url=http://upay.12580.com:80/upay/wps/service/tpfWapFormTrans.xhtml&wap_name=ChinaMoblie"/>
						<input type="hidden" id="we_return_url" name="WE_RETURN_URL" value="">
						<input type="hidden" name="BOM_TYP" value="wap"/>
						<input type="hidden" name="appid" value="" id="appid"/>
						<input type="hidden" id="isGzhzf" name="isGzhzf" value="0" /><!-- 是否为公众号支付，0为否，1为是 -->
				<input type="hidden" id="isXcxzf" name="isXcxzf" value="0" /><!-- 是否为小程序支付，0为否，1为是 -->
				</form>
				<form id="interSucForm" name="interSucForm"
						action="http://www.js.10086.cn/upay/wps/service/tpfInterfaceWapSerRes.xhtml" method="post" ><!-- accept-charset="utf-8" onsubmit="document.charset='utf-8';" -->
						<input type="hidden" name="ORD_NO" value="202104064729155538" /> 
						<input type="hidden" id="clickflg" name="CLICKFLG"/>
				</form>
				<input type="hidden" id="reserved1" value="" /><!--  -->
				<input type="hidden" id="CAP_CNL" value="WECHAT" />
				<input type="hidden" id="pag_ntf_url" value="https://wap.js.10086.cn/WSCZYLRESULT.thtml?ordno=F0BCB6E21DF03624F5F16E493F8A8E0A73A8756A170A726E3EBEB419B2D737BD&loginMobileAES=73A8756A170A726E3EBEB419B2D737BD&mobileAES=489EFF1F988A33E4D0B835273DF8027F&actualmoneyZFAES=06A0889555D932ED059C589C50C97E24&bankTypeAES=BB8C09BA0A84FE44EDFEE39687999590&ticketFlagAES=FF3B45660B219DD30C9459003540C9C4&bossDateAES=3BE13625B44A7A27A954C57F7D841066&isBindYxaAES=46411EE8A66F18B36B8A258402BE9749&isSHHYFlagAES=46411EE8A66F18B36B8A258402BE9749" /><!--  -->
				<input type="hidden" id="buttonType" value="1"/>
		
				<div class="modelBody" style="height:50%;top:50%;">
						<div style="width:100%;height: 270px;display:none;">
								<div style="height: 60px;"><!--  --></div>
								<div class="modelImg" style="width:70px;height:100px;"><img src="/upay/wps/images/modelIcon.png" alt="" width="70" height=""/></div>
								<div style="height: 20px;"><!--  --></div>
								<div style="text-align: center;font-size:16px;">支付平台跳转中，请继续支付</div>
						</div>
						<!-- <div style="height:25px;border-top:1px dotted #ccc"></div> -->
						<div style="text-align: center;">
								<a href="javascript:void(0)" class="ui_btn_model_g" onclick="finishPay()">我已完成支付</a>
								<span style="display: inline-block;width:10px;"></span>
								<a href="javascript:void(0)" class="ui_btn_model_g" onclick="giveupPay()">更换支付方式</a>
						</div>
				</div>
		</section>
		<div style="position:fixed;top:35%;left:0;width:100%;height:50px;">
				<img src="/upay/wps/images/loading_5g.gif" alt="" width="90" style="display:block;margin:0 auto;margin-bottom:15px;"/>
				<div style="text-align: center;font-size:14px;color:#999;">玩命加载中，请稍后...</div>
		</div>
		</body>
		</html>`,
	}
	doc, err := r.Document()
	if err != nil {
		panic(err)
	}
	form := r.ParseInputForm(doc, "payForm")
	fmt.Println(form.Stringify())
}
