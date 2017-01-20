$(function() {
	var baseUrl = 'http://v.prosnav.com/wxqyh/v1';

	page.base('/publish/newscontent');
	page('*', init);
	page('/', news);
	page();

	function init(ctx, next) {
	  ctx.query = qs.parse(location.search.slice(1));
	  next();
	}

	function news(ctx) {
		newsid = ctx.query.newsid;
		contentUrl = baseUrl + '/news/' + newsid;
		console.log(newsid);
		$.ajax({
			url: contentUrl,
			success: function(data, status, xhr) {
						$('#title').html(data.title);
						$('#content').find('span').html(data.crt.substr(0, 16).replace("T", " "));
						$('#content').find('div').html(data.content);
					},
			dataType: 'json',
			xhrFields: {
				withCredentials: true
			}
		});
	}

});
