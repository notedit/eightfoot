/* Foundation v2.2.1 http://foundation.zurb.com */
$(document).ready(function() {
	$('#login').click(function(e) {
		e.preventDefault();
		$('#myModal1').reveal();
	});
	$('#regist').click(function(e) {
		e.preventDefault();
		$('#myModal2').reveal();
	});
});