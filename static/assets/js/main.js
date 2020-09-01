!(function($) {
  "use strict";

  // Preloader
  $(window).on('load', function() {
    var preloader = $('#preloader');
    if (preloader.length) {
      preloader.delay(100).fadeOut('slow', function() {
        $(this).remove();
      });
    }
  });

  // Hero typed
  if ($('.typed').length) {
    var typed_strings = $(".typed").data('typed-items');
    typed_strings = typed_strings.split(',')
    new Typed('.typed', {
      strings: typed_strings,
      loop: true,
      typeSpeed: 100,
      backSpeed: 50,
      backDelay: 2000
    });
  }

  // Smooth scroll for the navigation menu and links with .scrollto classes
  $(document).on('click', '.nav-menu a, .scrollto', function(e) {
    if (location.pathname.replace(/^\//, '') == this.pathname.replace(/^\//, '') && location.hostname == this.hostname) {
      var target = $(this.hash);
      if (target.length) {
        e.preventDefault();

        var scrollto = target.offset().top;

        $('html, body').animate({
          scrollTop: scrollto
        }, 1500, 'easeInOutExpo');

        if ($(this).parents('.nav-menu, .mobile-nav').length) {
          $('.nav-menu .active, .mobile-nav .active').removeClass('active');
          $(this).closest('li').addClass('active');
        }

        var body = $('body');
        if (body.hasClass('mobile-nav-active')) {
          body.removeClass('mobile-nav-active');
          $('.mobile-nav-toggle i').toggleClass('icofont-navigation-menu icofont-close');
        }
        return false;
      }
    }
  });

  // Activate smooth scroll on page load with hash links in the url
  $(document).ready(function() {
    if (window.location.hash) {
      var initial_nav = window.location.hash;
      if ($(initial_nav).length) {
        var scrollto = $(initial_nav).offset().top;
        $('html, body').animate({
          scrollTop: scrollto
        }, 1500, 'easeInOutExpo');
      }
    }
  });

  $(document).on('click', '.mobile-nav-toggle', function(e) {
    $('body').toggleClass('mobile-nav-active');
    $('.mobile-nav-toggle i').toggleClass('icofont-navigation-menu icofont-close');
  });

  $(document).click(function(e) {
    var container = $(".mobile-nav-toggle");
    if (!container.is(e.target) && container.has(e.target).length === 0) {
      var body = $('body');
      if (body.hasClass('mobile-nav-active')) {
        body.removeClass('mobile-nav-active');
        $('.mobile-nav-toggle i').toggleClass('icofont-navigation-menu icofont-close');
      }
    }
  });

  // Navigation active state on scroll
  var nav_sections = $('section');
  var main_nav = $('.nav-menu, #mobile-nav');

  $(window).on('scroll', function() {
    var cur_pos = $(this).scrollTop() + 300;

    nav_sections.each(function() {
      var top = $(this).offset().top,
        bottom = top + $(this).outerHeight();

      if (cur_pos >= top && cur_pos <= bottom) {
        if (cur_pos <= bottom) {
          main_nav.find('li').removeClass('active');
        }
        main_nav.find('a[href="#' + $(this).attr('id') + '"]').parent('li').addClass('active');
      }
      if (cur_pos < 200) {
        $(".nav-menu ul:first li:first").addClass('active');
      }
    });
  });

  // Back to top button
  $(window).scroll(function() {
    if ($(this).scrollTop() > 100) {
      $('.back-to-top').fadeIn('slow');
    } else {
      $('.back-to-top').fadeOut('slow');
    }
  });

  $('.back-to-top').click(function() {
    $('html, body').animate({
      scrollTop: 0
    }, 1500, 'easeInOutExpo');
    return false;
  });

  // jQuery counterUp
  $('[data-toggle="counter-up"]').counterUp({
    delay: 10,
    time: 1000
  });

  // Skills section
  $('.skills-content').waypoint(function() {
    $('.progress .progress-bar').each(function() {
      $(this).css("width", $(this).attr("aria-valuenow") + '%');
    });
  }, {
    offset: '80%'
  });

  // Init AOS
  function aos_init() {
    AOS.init({
      duration: 1000,
      once: true
    });
  }

  // Porfolio isotope and filter
  $(window).on('load', function() {
    var portfolioIsotope = $('.portfolio-container').isotope({
      itemSelector: '.portfolio-item'
    });

    $('#portfolio-flters li').on('click', function() {
      $("#portfolio-flters li").removeClass('filter-active');
      $(this).addClass('filter-active');

      portfolioIsotope.isotope({
        filter: $(this).data('filter')
      });
      aos_init();
    });

    // Initiate venobox (lightbox feature used in portofilo)
    $('.venobox').venobox({
      'share': false
    });

    // Initiate aos_init() function
    aos_init();

  });

  // Testimonials carousel (uses the Owl Carousel library)
  $(".testimonials-carousel").owlCarousel({
    autoplay: true,
    dots: true,
    loop: true,
    items: 1
  });

  // Portfolio details carousel
  $(".portfolio-details-carousel").owlCarousel({
    autoplay: true,
    dots: true,
    loop: true,
    items: 1
  });

})(jQuery);

function addNewGameResult() {
  var row = $('.game-result-row:last');
  var newRow = row.clone(),
      rowParent = row.parent(),
      validateRow = row.next('.validate').clone();

  newRow.appendTo(rowParent);
  newRow.hide().slideDown();
  validateRow.html('');
  validateRow.appendTo(rowParent);
  newRow.find('.amount-group input').val(null);

  newRow.find('.delete-row-btn').css('pointer-events', 'auto');
  newRow.find('.delete-row-btn').css('opacity', 'inherit');
}

  function onRowDelete(event) {
    var currentRow = $(event.currentTarget).closest('.game-result-row');
    var validateRow = currentRow.next('.validate');

    validateRow.hide('blind', function () {
      validateRow.remove();
    });
    currentRow.hide('blind', function f() {
      currentRow.remove();
    });

    event.preventDefault();
  }

  function validateGameName(gameName, validationBlock, errorMessage) {
    if(!gameName || gameName.length < 5) {
      validationBlock.html(errorMessage).slideDown();
      return false;
    } else {
      validationBlock.hide('blind');
      return true;
    }
  }

  function validateGameDate(gameDate, validationBlock, errorMessage) {
    if(!gameDate || gameDate.length === 0) {
      validationBlock.html(errorMessage).slideDown();
      return false;
    } else {
      validationBlock.hide('blind');
      return true;
    }
  }

  function validatePlayer (playerId, validationBlock, errorMessage) {
    // noinspection EqualityComparisonWithCoercionJS
    if (!playerId || playerId == 0) {
      validationBlock.html(errorMessage).slideDown();
      return false;
    } else {
      validationBlock.hide('blind');
      return true;
    }
  }

  function validateResultRow(winnerId, looserId, validationBlock, errorMessage) {
    if(winnerId >0 && looserId > 0 && winnerId === looserId) {
      validationBlock.html(errorMessage).slideDown();
      return false;
    } else {
      validationBlock.hide('blind');
      return true;
    }
  }

  function validateAmount(amount, validationBlock, errorMessage, negativeValueErrorMessage) {
    if (!amount) {
      validationBlock.html(errorMessage).slideDown();
      return false;
    } else if(amount <= 0) {
      validationBlock.html(negativeValueErrorMessage).slideDown();
    } else {
      validationBlock.hide('blind');
      return true;
    }
  }

  function validateAllResults(gameResults, validationBlock) {
    for(var i = 0; i < gameResults.length; i++) {
      var iRow = gameResults[i];

      for (var j = i + 1; j < gameResults.length; j ++){
        var wrongPlayerId, jRow = gameResults[j];

        if(iRow.winnerId === jRow.looserId){
          wrongPlayerId = iRow.winnerId;
        } else if (iRow.looserId === jRow.winnerId) {
          wrongPlayerId = iRow.looserId;
        } else {
          continue;
        }

        var wrongPlayerName = $('select[name=winner]:first').find('option[value="' + wrongPlayerId + '"]').text();
        validationBlock.html('Игрок ' + wrongPlayerName + ' не может одновременно победить и проиграть').slideDown();
        return false;
      }
    }

    validationBlock.hide('blind');
    return true;
  }

  function onNewGameResultSubmit(event) {
    event.preventDefault();

    var gameName = $('#name'),
        gameDate = $('#game-date'),
        this_form = $(this),
        gameResults = [],
        isValidationOk = true;

    // noinspection JSBitwiseOperatorUsage
    isValidationOk = validateGameName(gameName.val(), gameName.next('.validate'), gameName.attr('data-msg')) & validateGameDate(gameDate.val(), gameDate.parent().next('.validate'), gameDate.attr('data-msg'));

    $('#game-result-form .game-result-row').each(function (index, item) {
      var winnerId, looserId, amount, isRowValidationOk,
          winnerElement = $(item).find('select[name=winner]'),
          looserElement = $(item).find('select[name=looser]'),
          amountElement = $(item).find('input[name=amount]');

      winnerId = parseInt(winnerElement.find('option').filter(':selected').val());
      isRowValidationOk = validatePlayer(winnerId, winnerElement.next('.validate'), winnerElement.attr('data-msg'));

      looserId = parseInt(looserElement.find('option').filter(':selected').val());
      isRowValidationOk = validatePlayer(looserId, looserElement.next('.validate'), looserElement.attr('data-msg')) && isRowValidationOk;

      amount = parseInt(amountElement.val());
      isRowValidationOk = validateAmount(amount, amountElement.next('.validate'), amountElement.attr('data-msg'), 'Выигрыш должен быть больше 0') && isRowValidationOk;

      if (isRowValidationOk) {
        isRowValidationOk = validateResultRow(winnerId, looserId, $(item).next('.validate'), 'Сам у себя победил?') && isRowValidationOk;
      }

      isValidationOk = isValidationOk && isRowValidationOk;
      gameResults.push({winnerId: winnerId, looserId: looserId, amount: amount});
    });

    isValidationOk = validateAllResults(gameResults, $('#common-validation')) && isValidationOk;

    if (isValidationOk) {
      var formData = { gameName: gameName.val(), gameDate: gameDate.val(), gameResults: JSON.stringify(gameResults) }
      addGameResult(this_form.attr('action'), JSON.stringify(formData), this_form)
    }
  }

  function addGameResult(action, data, this_form) {
    var submitButton = $('#game-result-form button[type=submit]');
    submitButton.prop('disabled', true);
    submitButton.css('opacity', '0.5');

    $.ajax({
      type: "POST",
      url: action,
      data: data,
      contentType: "application/x-www-form-urlencoded; charset = UTF-8",
      success: function(data, textStatus ){
        console.log("success data");
        console.log(textStatus);
        if (textStatus === 'success') {
          location.reload();
          // this_form.find('.loading').slideUp();
          // this_form.find('.sent-message').slideDown();
          // this_form.find("input:not(input[type=submit]), textarea").val('');
        } else {
          this_form.find('.loading').slideUp();
          this_form.find('.error-message').slideDown().html(textStatus);
        }
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
      error:function (data, a, b, c, d) {
        console.log(data);
        console.log(a);
        console.log(b);
        console.log(c);
        console.log(d);
      }
    });
  }