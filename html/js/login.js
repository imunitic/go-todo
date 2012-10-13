(function() {
	$(function() {
		$("#loginFormErrors").toggle();
		$("#login").submit(function(e) {
			e.preventDefault();
			var data = $(this).serialize();
			$.post($(this).attr("action"), data, function(json) {
				if(typeof(json.error) != 'undefined') {					
					$("#loginFormErrors").toggle();
					$("#loginFormErrors").html(json.error);
				}
			})
		});
	});
})()