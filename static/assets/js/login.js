function onSigningIn(event) {
    event.preventDefault();
    var login = $("#login").val()
    var password = $("#password").val()
    var loginButton = $("#login_button")

    loginButton.prop('disabled', true);
    loginButton.css('opacity', '0.5');

    var data = JSON.stringify({login: login, password: password})
    console.log(data)

    $.ajax({
        type: "POST",
        url: "login",
        data: data,
        contentType: "application/x-www-form-urlencoded; charset = UTF-8",
        success: function(){
            location.assign('/')
        },
        done: function(data) {
            console.log('done data')
            console.log(data)
            location.reload();
            var error_msg = "Form submission failed!<br>";
            if(data.statusText || data.status) {
                error_msg += 'Status:';
                if(data.statusText) {
                    error_msg += ' ' + data.statusText;
                }
                if(data.status) {
                    error_msg += ' ' + data.status;
                }
                error_msg += '<br>';
            }
            if(data.responseText) {
                error_msg += data.responseText;
            }
            this_form.find('.loading').slideUp();
            this_form.find('.error-message').slideDown().html(error_msg);
        },
        error:function (data) {
            if (data.status === 401) {
                loginButton.prop('disabled', false);
                loginButton.css('opacity', '1');
                var validationMessageBlock = $('#validation-message');
                validationMessageBlock.slideDown();
                setTimeout(function(){validationMessageBlock.toggle('blind')}, 10);
            }
        }
    });
}

function onLoginKeyDown(e) {
    if (e.which === 32)
        return false;
}

function onLoginChange() {
    var loginInput = $("input#login");
    loginInput.val(loginInput.val().replace(/\s/g, ""));
}