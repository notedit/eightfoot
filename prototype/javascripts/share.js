/* Foundation v2.2.1 http://foundation.zurb.com */

$(document).ready(function() {
	$('#myTags').on( 'click', 'a', function(e){
		e.preventDefault();
        var $this = $(this),
        	value = $this.attr( 'title' );
        
		if ( !$this.is( '.added' ) ) {
			var tags = $('#shareTags').val() + ' "' + value + '" ';
			$this.removeClass('blue');
			$this.addClass('added white');
			$('#shareTags').val(tags);
		} else {
			$this.removeClass('added white');
			$this.addClass('blue');
			tags = ;
		}
	});
/*
	var $tags = $('#myTags');
	$tags.on('click','a',function(e) {
                
                var $this = $(this),
                value = $this.text(),
                id = $this.data('id');
                
                if ($this.is('.badge-success')) {
                    $('#tagid').remove();
                    $this.removeClass('badge-success');
                } else {
                    $('#tagid').remove();
                    $tags.append(format(tagHtml,{tid:id}));
                    $this.addClass('badge-success');
                    console.log('badge-success');
                    $('a',$tags).each(function(i){
                        $that = $(this);
                        if ($that.is('.badge-success') && $this[0] !== $that[0]) {
                            $that.removeClass('badge-success');
                        }
                    });
                }
            });*/
});