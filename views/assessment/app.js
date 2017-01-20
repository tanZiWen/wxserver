$(function () {
    var host = 'http://121.40.88.18:4001/wxqyh/v1/',
        assessmentApi = host + 'assessment',
        paperApi = assessmentApi + '/paper'

    $("#xyLi").click(function () {
        $("#xypopDiv").show();
    });
    $("#ikonwImg").click(function () {
        if ($("input[name='xy_checkbox']:checked").length != 3) {
            alert("请勾选并同意利得承诺函选项！");
        } else {
            $("#xypopDiv").hide();
        }
    });
    function getPageData() {
        var data = {}
        data.q1 = $('input:radio[name="q1"]:checked').val();
        data.q2 = $('input:radio[name="q2"]:checked').val();
        data.q3 = $('input:radio[name="q3"]:checked').val();
        data.q4 = $('input:radio[name="q4"]:checked').val();
        data.q5 = $('input:radio[name="q5"]:checked').val();
        data.q6 = $('input:radio[name="q6"]:checked').val();
        data.q7 = $('input:radio[name="q7"]:checked').val();
        data.q8 = $('input:radio[name="q8"]:checked').val();
        data.q9 = $('input:radio[name="q9"]:checked').val();
        data.q10 = $('input:radio[name="q10"]:checked').val();
        data.q11 = $('input:radio[name="q11"]:checked').val();
        data.q12 = $('input:radio[name="q12"]:checked').val();
        data.q13 = $('input:radio[name="q13"]:checked').val();
        data.q14 = $('input:radio[name="q14"]:checked').val();
        data.q15 = $('input:radio[name="q15"]:checked').val();
        data.q16 = $('input:radio[name="q16"]:checked').val();
        data.q17 = $('input:radio[name="q17"]:checked').val();
        data.q18 = $('input:radio[name="q18"]:checked').val();
        data.mobile = $("#mobile").val();
        data.name = $('#name').val();
        data.email = $('#email').val();
        return data;
    }

    function parseQuestion(seq, name, question) {
        var qt = '<ul>';
        qt += seq + '.' + ' ' + question.QUESTION + '<br><br>';
        var i = 0;
        for (attr in question) {
            if (attr == 'QUESTION') {
                continue;
            }
            optionId = name + '_' + attr
            qt += '<li><input id= "' + optionId + '"" type="radio" style="width:18px;height:18px;" name="' + name + '" value="' + attr + '"><label for="' + optionId + '">' + attr + ' ' + question[attr] + '</label></li>'
        }
        qt += '<br><br>';
        return qt
    }

    function initPaper(data) {
        var paper = "", i = 1;
        for (q in data) {
            sec = 'q' + i;
            paper += parseQuestion(i, sec, data[sec])
            i++;
        }
        $('#paper').html(paper);
    }

    var sumbitFlag = false;

    function init() {
        $('.container').hide();
        $('#mobile').val('');
        $('#name').val('');
        $.ajax({
            url: paperApi,
            method: 'GET',
            dataType: 'json',
            xhrFields: {
                withCredentials: true
            },
            success: function (data, status, jqXHR) {
                initPaper(data);
            }
        });
        $('#main').show();
        sumbitFlag = false;
    }

    $('#submit').click(function () {
        if (validate() == false) {
            return false;
        }
        // 避免重复提交。
        if (sumbitFlag == true) {
            return false;
        }
        $.ajax({
            url: assessmentApi,
            method: 'POST',
            dataType: 'json',
            xhrFields: {
                withCredentials: true
            },
            data: getPageData(),
            success: function (data, status, jqXHR) {
                $('#type').html("您的投资类别为："+data.type);
                $('#score').html("您的测评总分为："+data.score+"分");
                $('#endurance').html("您的风险承受能力为："+data.endurance);
                $('#venture').html("适合您的基金产品类型为："+data.venture);
                $('.container').hide();
                $('#success').show();
                sumbitFlag == true;
            }
        });
    });
    $('#reassess').click(function () {
        init();
    });
    function validateEmail(val) {
        var reg = /^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+\.[a-zA-Z]{2,3}$/;
        if(!reg.test(val)) {
            return false;
        }
        return true;
    }
    function validate() {
        var ask_val_1 = $('input:radio[name="q1"]:checked').val();
        var ask_val_2 = $('input:radio[name="q2"]:checked').val();
        var ask_val_3 = $('input:radio[name="q3"]:checked').val();
        var ask_val_4 = $('input:radio[name="q4"]:checked').val();
        var ask_val_5 = $('input:radio[name="q5"]:checked').val();
        var ask_val_6 = $('input:radio[name="q6"]:checked').val();
        var ask_val_7 = $('input:radio[name="q7"]:checked').val();
        var ask_val_8 = $('input:radio[name="q8"]:checked').val();
        var ask_val_9 = $('input:radio[name="q9"]:checked').val();
        var ask_val_10 = $('input:radio[name="q10"]:checked').val();
        var ask_val_11 = $('input:radio[name="q11"]:checked').val();
        var ask_val_12 = $('input:radio[name="q12"]:checked').val();
        var ask_val_13 = $('input:radio[name="q13"]:checked').val();
        var ask_val_14 = $('input:radio[name="q14"]:checked').val();
        var ask_val_15 = $('input:radio[name="q15"]:checked').val();
        var ask_val_16 = $('input:radio[name="q16"]:checked').val();
        var ask_val_17 = $('input:radio[name="q17"]:checked').val();
        var ask_val_18 = $('input:radio[name="q18"]:checked').val();
        var mobile = $('#mobile').val();
        var name = $('#name').val();
        var email = $('#email').val();

        if (ask_val_1 == null) {
            alert("请选择题目1的答案");
            return false;
        }
        if (ask_val_2 == null) {
            alert("请选择题目2的答案");
            return false;
        }
        if (ask_val_3 == null) {
            alert("请选择题目3的答案");
            return false;
        }
        if (ask_val_4 == null) {
            alert("请选择题目4的答案");
            return false;
        }
        if (ask_val_5 == null) {
            alert("请选择题目5的答案");
            return false;
        }
        if (ask_val_6 == null) {
            alert("请选择题目6的答案");
            return false;
        }

        if (ask_val_7 == null) {
            alert("请选择题目7的答案");
            return false;
        }

        if (ask_val_8 == null) {
            alert("请选择题目8的答案");
            return false;
        }

        if (ask_val_9 == null) {
            alert("请选择题目9的答案");
            return false;
        }

        if (ask_val_10 == null) {
            alert("请选择题目10的答案");
            return false;
        }

        if (ask_val_11 == null) {
            alert("请选择题目11的答案");
            return false;
        }
        if (ask_val_12 == null) {
            alert("请选择题目12的答案");
            return false;
        }
        if (ask_val_13 == null) {
            alert("请选择题目13的答案");
            return false;
        }
        if (ask_val_14 == null) {
            alert("请选择题目14的答案");
            return false;
        }
        if (ask_val_15 == null) {
            alert("请选择题目15的答案");
            return false;
        }
        if (ask_val_16 == null) {
            alert("请选择题目16的答案");
            return false;
        }
        if (ask_val_17 == null) {
            alert("请选择题目17的答案");
            return false;
        }
        if (ask_val_18 == null) {
            alert("请选择题目18的答案");
            return false;
        }
        if (!mobile) {
            alert("请输入您的手机号码");
            return false;
        }
        if (!name) {
            alert("请输入您的姓名");
            return false;
        }
        if(!validateEmail(email)) {
            alert("请输入正确的邮箱");
            return false;
        }

    }

    init();
});