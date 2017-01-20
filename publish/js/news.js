$(function(){
	baseUrl = 'http://v.prosnav.com/wxqyh/v1'

	function createNewsItem(field) {
		liEle = $('<li></li>');
		liEle.addClass('ui-first-child').addClass('ui-last-child');
		aEle = $('<a href="#newsContent" data-transition="slide"></a>');
		aEle.addClass('ui-btn').addClass('ui-btn-icon-right').addClass('ui-icon-carat-r');
		h2Ele = $('<h2>'+ field.title +'</h2>');
		pEle = $('<p>'+ field.summary +'</p>');
		pEle.addClass('summary');
		aEle.append(h2Ele);
		aEle.append(pEle);
		liEle.append(aEle);
		
		aEle.click(function() {
			loadNewsContent(field.id)
		});

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

	function createNewsForm(currPage, pageCount) {
		return {"currPage": currPage,
				 "pageCount": pageCount
				};
	}

	function DataInfo(currPage, pageCount, total) {
		this.currPage = currPage;
		this.pageCount = pageCount;
		this.total = total || 0;
	}

	function appendNews(newsObject) {
		dataLog.total = newsObject.total;
		$('#username').html(newsObject.username);
		if (newsObject.news) {
		    $.each(newsObject.news, function(i, field){
	                $("#news").append(createNewsItem(field));
	            });
		}
	}

	var dataLog = new DataInfo(1, 5);
	dataLog.isLastPage = function() {
		if (this.currPage * this.pageCount >= this.total) {
			return true;
		}
		return false;
	};

	function loadData() {
		$.ajax({
			url: baseUrl + '/news/',
			data: createNewsForm(dataLog.currPage, dataLog.pageCount),
			success: function(data, status, xhr) {
						console.log(data);
						appendNews(data);
						if (!data.news) {
							$('#empty').show();
							return
						}
						$('#more').show();
						dataLog.total = data.total;
						if (dataLog.isLastPage()) {
							$('#more').remove();
						}
					},
			xhrFields: {
				withCredentials: true
			},
			dataType: 'json',
			error: function(jqxhr, status, err) {
				alert("Internal error");
			}
		});
	}

	loadData();

	$('#more').click(function () {
			dataLog.currPage++;
			loadData();
	});

});
