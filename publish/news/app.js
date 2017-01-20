$(function(){

	var baseUrl = 'http://v.prosnav.com/wxqyh/v1',
		perpage = 5,
		total = 0,
		currpage = 1;

	page.base('/publish/news');
	page('*', init);
	page();
	newslist();

	function init(ctx, next) {
		ctx.query = qs.parse(location.search.slice(1));
		if (Object.keys(ctx.query).length) {
			currpage = ~~ctx.query.page;
		}
		next();
	}

	function newslist() {
		adjustPager();
		$.ajax({
			url: baseUrl + '/news/',
			data: {"currPage": currpage, "pageCount": perpage},
			success: function(data, status, xhr) {
				show(data);
			},
			dataType: 'json',
			xhrFields: {
				withCredentials: true
			},
			error: function(jqxhr, status, err) {
				alert("Internal error");
			}
		});
	}



	function adjustPager() {
		$('#prev').attr('href', '/publish/news?page=' + (currpage - 1));
		$('#next').attr('href', '/publish/news?page=' + (currpage + 1));
	}

	function show(data) {
		$("#news").html('');
		$('#username').html(data.username+',您好');
		if (data.news) {
			total = data.total;
		    $.each(data.news, function(i, field){
	            $("#news").append(createNewsItem(field));
	        });
		}
		if (currpage == 1) {
			$('#prev').hide();
		} else {
			$('#prev').show();
		}

		if (currpage * perpage >= total) {
			$('#next').hide();
		} else {
			$('#next').show();
		}
	}

	function createNewsItem(data) {
		h2Ele = $('<h2>' + data.title + '</h2>');
		pEle = $('<p>'+ data.summary +'</p>').addClass('summary');
		aEle = $('<a href="/publish/newscontent/?newsid='+data.id+'"></a>').addClass('ui-btn').append(h2Ele).append(pEle);
		liEle = $('<li></li>').append(aEle).addClass('ui-first-child').addClass('ui-last-child').addClass('ui-btn');

		return liEle;
	}

	function loadNewsContent(newsid) {
		contentUrl = baseUrl + '/news/' + newsid;
		$.ajax({
			url: contentUrl,
			success: function(data, status, xhr) {
						$('#newsContent').find('#title').html(data.title);
						$('#newsContent').find('span').html(data.crt.substr(0, 16).replace("T", " "));
						$('#newsContent').find('#content').find('div').html(data.content);
						$('#aniimated-thumbnials').find('a').addClass("imageGap");
						$('#aniimated-thumbnials').lightGallery({
							thumbnail:true,
							animateThumb: false,
							showThumbByDefault: false,
							download: false,
							zoom: true,
							controls: false,
							hideBarsDelay:1000
						});
					},
			dataType: 'json',
			xhrFields: {
				withCredentials: true
			}
		});
	}
});
