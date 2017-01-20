$(function(){

	$.ajaxSetup({
		error:function(qXHR, exception) {
			if (qXHR.status == 401 || qXHR.statusText == "error") {
				window.location = '/publish/login.html';
			}
		}
	});

	baseUrl = 'http://v.prosnav.com/wxqyh/v1'

	var ckeditor = CKEDITOR.replace( 'ckeditor1' );
	function DataInfo(currPage, pageCount, total) {
		this.currPage = currPage;
		this.pageCount = pageCount;
		this.total = total || 0;
	}
	
	var dataLog = new DataInfo(1, 10);
	dataLog.isFirstPage = function() {
		return this.currPage == 1;
	};
	dataLog.isLastPage = function() {
		if (this.currPage * this.pageCount >= this.total) {
			return true;
		}
		return false;
	};

	function createNewsForm(currPage, pageCount) {
		return {"currPage": currPage,
				 "pageCount": pageCount
				};
	}

	function loadData() {
		$.ajax({
			url: baseUrl + '/news/',
			type: 'POST',
			data: createNewsForm(dataLog.currPage, dataLog.pageCount),
			success: function(data, status, xhr) {
						console.log(data);
						appendNews(data);
						dataLog.total = data.total;
					},
			dataType: 'json',
			xhrFields: {
				withCredentials: true
			}
		});
	}

	function appendNews(newsObject) {
		dataLog.total = newsObject.total;
		if (newsObject.news) {
			$("#news").html('')
		    $.each(newsObject.news, function(i, field){
	            $("#news").append(createNewsItem(field));
	        });
		}
	}

	function createNewsItem(field) {
		row = $('<tr></tr>');
		title = $('<td>'+ field.title +'</td>');
		summary = $('<td>'+ field.summary +'</td>');
		pcode = $('<td>'+ field.productcode +'</td>');
		operation = $('<td></td>');
		edit = $('<i class="icon-edit" data-toggle="modal" data-target=".newsmodal"></i>');
		minus = $('<i class="icon-cancel"></i>');

		edit.click(function() {
			loadNewsContent(field.id);
		});

		minus.click(function() {
			bootbox.confirm({ 
			    size: 'small',
			    message: "Are you sure?", 
			    className: 'confirm',
			    callback: function(result){ /* your callback code */ 
			    	if (result) {
						deleteNews(field.id);
					}
			    }
			})
		});

		operation.append(edit).append(minus);
		row.append(title).append(summary).append(pcode).append(operation)

		return row;
	}
	$('#add').click(function() {
		cleanModal();
	});

	function cleanModal() {
		$('#newsid').val('');
		$('#newstitle').val('');
		$('#newssummary').val('');
		$('#newspcode').val('');
		ckeditor.setData("");
	}

	function deleteNews(newsid) {
		contentUrl = baseUrl + '/news/delete/' + newsid;
		$.ajax({
			url: contentUrl,
			type:'POST',
			success: function(data, status, xhr) {
						dataLog.currPage = 1;
						loadData();
					},
			dataType: 'json',
			xhrFields: {
				withCredentials: true
			}
		});
	}

	function loadNewsContent(newsid) {
		contentUrl = baseUrl + '/news/' + newsid;
		$.ajax({
			url: contentUrl,
			success: function(data, status, xhr) {
						$('#newsid').val(newsid);
						$('#newstitle').val(data.title);
						$('#newssummary').val(data.summary);
						$('#newspcode').val(data.productcode);
						ckeditor.setData(data.content);
						console.log(data);
					},
			dataType: 'json',
			xhrFields: {
				withCredentials: true
			}
		});
	}

	function addNews() {
		contentUrl = baseUrl + '/news/news';
		$.ajax({
			url: contentUrl,
			type: 'POST',
			data:{
					"title": $('#newstitle').val(),
					"summary": $('#newssummary').val(),
					"productcode": $('#newspcode').val(),
					"content":ckeditor.getData()
				},
			success: function(data, status, xhr) {
						dataLog.currPage = 1;
						loadData();
					},
			dataType: 'json',
			xhrFields: {
				withCredentials: true
			}
		});
	}

	function modifyNews() {
		contentUrl = baseUrl + '/news/update/' + $('#newsid').val();
		$.ajax({
			url: contentUrl,
			type: 'POST',
			data:{
					"title": $('#newstitle').val(),
					"summary": $('#newssummary').val(),
					"productcode": $('#newspcode').val(),
					"content":ckeditor.getData()
				},
			success: function(data, status, xhr) {
						dataLog.currPage = 1;
						loadData();
					},
			dataType: 'json',
			xhrFields: {
				withCredentials: true
			}
		});
	}

	$('#save').click(function() {
		if ($('#newsid').val()) {
			modifyNews();
		} else {
			addNews();
		}
		$('.modal').modal('hide');
	});
	
	$('#pre').tooltip({
		title:'First Page',
		placement:'top',
		trigger: 'manual'
	});

	$('#next').tooltip({
		title:'Last Page',
		placement:'bottom',
		trigger: 'manual'
	});

	$('#pre').click(function() {
		if (dataLog.isFirstPage()) {
			var obj = $(this)
			obj.tooltip('show');
			//$(this).tooltip('hide');
			setTimeout(function() {
				obj.tooltip('hide');
			}, 500);
			return;
		}
		dataLog.currPage--;
		loadData();
	});

	$('#next').click(function() {
		if (dataLog.isLastPage()) {
			var obj = $(this)
			obj.tooltip('show');
			setTimeout(function() {
				obj.tooltip('hide');
			}, 500);
			return;
		}
		dataLog.currPage++;
		loadData();
	})

	
	loadData();

});
