$(function(){
	var host = 'http://v.prosnav.com/wxqyh/v1/act',
	initDataApi = host + '/meta',
	queryActivityApi = host + '/activity',
	signInApi = host + '/signin';

	$(document).ajaxError(function() {
	  	$('.container').hide();
		$('#error').show();
		closeWindow();
	});

	function onMetaDataResponse(data) {
		wx.config({
		    debug: false, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
		    appId: data.appid, // 必填，企业号的唯一标识，此处填写企业号corpid
		    timestamp: data.timestamp, // 必填，生成签名的时间戳
		    nonceStr: data.nonceStr, // 必填，生成签名的随机串
		    signature: data.signature,// 必填，签名，见附录1
		    jsApiList: ['getLocation'] // 必填，需要使用的JS接口列表，所有JS接口列表见附录2
		});

		wx.ready(function(){
		    wx.getLocation({
			    type: 'wgs84', // 默认为wgs84的gps坐标，如果要返回直接给openLocation用的火星坐标，可传入'gcj02'
			    success: function(res) {
			        $.ajax({
						url: queryActivityApi,
						method: 'GET',
						dataType: 'json',
						data: {latitude: res.latitude, longitude: res.longitude}, 
						xhrFields: {
							withCredentials: true
						},
						success: function(ret, status, jqXHR) {
							if (ret.signed) {
								$('.container').hide();
								$('#signed').show();
								closeWindow();
								return;
							}
				        	if (ret.act) {
				        		$('#activityName').html(ret.act.activityname);
				        		$('#actid').val(ret.act.id);
				        		$('#custName').val(ret.username);
				        		$('#mobile').val(ret.mobile);
				        		$('.container').hide();
								$('#main').show();
				        		return;
				        	}
				        	$('.container').hide();
				        	$('#noActivity').show();
				        	closeWindow();
				        }
			        });
			    }
			});
		});
	};

	function closeWindow() {
		setTimeout(function() {
				WeixinJSBridge.call('closeWindow');
		}, 1000);
	}

	function start() {
		$('.container').hide();
		$('#loading').show();
		$.ajax({
			url: initDataApi+'?url='+window.location,
			dataType: 'json',
			method: 'GET',
			xhrFields: {
				withCredentials: true
			},
			success: onMetaDataResponse
		});
	}
	$(":text").blur(function() {
		if (!$(this).val()) {
			$(this).focus();
		}
	}); 
	$('#submit').click(function() {
		$.ajax({
			url: signInApi, 
			method: 'POST',
			dataType: 'json',
			xhrFields: {
				withCredentials: true
			},
			data: { activityid: $('#actid').val(), custname: $('#custName').val(), mobile: $('#mobile').val() }, 
			success: function(data, status, jqXHR) {
				if (data == 'signed_in') {
					$('.container').hide();
					$('#signed').show();
					closeWindow();
				} else {
					$('.container').hide();
					$('#success').show();
					closeWindow();
				}
			}
		});
	});
	start();
});