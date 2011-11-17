$(document).ready(function(){
   // Remove last three elements
   $('#friends').children().each(function(key, element){
       if(key > 0) {
           $(this).remove();
       }
   })

   // Add the plus sign at the end of the first element
   $('#friends li').append('<span class="toggle add">+</span>');

   // Add inputs if needed
   var inputs = 1;
   $('.add').live('click', function(){
      if(inputs <= 9) {
          $(this).removeClass('add');
          $(this).addClass('del');
          $(this).text('x');
          $('#friends').append('<li><label>Email Friend</label><input name="friends[]" type="email"><span class="toggle add">+</span></li>');
          inputs++;
    } else {
        alert("Don't be silly, nobody has so many friends");
    }
   });
   
   // Delete the selected input
   $('.del').live('click', function(){
       $(this).parent().remove();
   });
   
   
});