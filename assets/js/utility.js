function stripDecimals(t){return 0|t}function customSessionRequest(t){var e=localStorage.getItem("city"),a=localStorage.getItem("country"),o=localStorage.getItem("battery"),n=localStorage.getItem("screen"),c=localStorage.getItem("page"),r=localStorage.getItem("ip");$.ajax({async:!0,type:"POST",dataType:"json",url:"/"+t,data:{info:"vsession",city:e,country:a,screen_size:n,page:c,battery:o,ip:r},success:function(){}}).done(function(t){})}function MoneyFormat(t){return"number"==typeof t||"string"==typeof t&&(t=parseInt(t,10)),t.toFixed(2).replace(/\d(?=(\d{3})+\.)/g,"$&,")}function updateCart(t){var e=t.length,a=$(".cart-dropdown");$(".count-label").text(e);for(var o="",n=0,c=0;c<t.length;c++){var r=t[c].product_name,s=t[c].product_image,i=t[c].doc_number,l=t[c].product_id,g=t[c].quantity,u=t[c].payable_amount,d=t[c].trxid,m=Number(u)/Number(g);n+=Number(u),o+='<div class="entry">'+`<div class="entry-thumb"><a href="#"><img src="${s}" alt="Product"></a></div>`+'<div class="entry-content">'+`<h4 class="entry-title"><a href="#">${r}</a></h4><span class="entry-meta">${g} x ৳${MoneyFormat(m)}</span>`+"</div>"+`<div class="entry-delete" data-tid="${d}" data-doc="${i}" data-item="${l}"><i class="icon-x"></i></div>`+"</div>\n\n"}o+='<div class="text-right">'+`<p class="text-gray-dark py-2 mb-0"><span class='text-muted'>Subtotal:</span> &nbsp;৳${MoneyFormat(n)}</p>`+"</div>",o+='<div class="d-flex"><div class="pr-2 w-50"><a class="btn btn-secondary btn-sm btn-block mb-0" href="/cart">Expand Cart</a></div><div class="pl-2 w-50"><a class="btn btn-primary btn-sm btn-block mb-0" href="/checkout">Checkout</a></div></div>',a.empty().append(o)}function clearData(){localStorage.clear()}function screenWidthHeight(){var t=window.screen.height,e=`${window.screen.width}x${t}`;return localStorage.setItem("screen",e),e}function currentUrl(){var t=window.location.pathname;return localStorage.setItem("page",t),t}function checkIfEmpty(){var t=localStorage.getItem("city"),e=localStorage.getItem("ip");localStorage.getItem("country");try{0==t.length||0==e.length?(getBatteryStatus(),screenWidthHeight(),currentUrl(),getLocation()):sessionRequest()}catch(t){getBatteryStatus(),screenWidthHeight(),currentUrl(),getLocation()}}function getBatteryStatus(){var t=!1;return navigator.getBattery&&(t=!0,navigator.getBattery().then(function(t){})),localStorage.setItem("battery",t),t}function getLocation(){$.ajax({url:"https://geolocation-db.com/jsonp",jsonpCallback:"callback",dataType:"jsonp",success:function(t){}}).done(function(t){city=t.city,country_name=t.country_name,localStorage.setItem("city",city),localStorage.setItem("country",country_name),localStorage.setItem("lat",t.latitude),localStorage.setItem("long",t.longitude),localStorage.setItem("ip",t.IPv4),sessionRequest()})}function sessionRequest(){var t=localStorage.getItem("city"),e=localStorage.getItem("country"),a=localStorage.getItem("battery"),o=localStorage.getItem("screen"),n=localStorage.getItem("page"),c=localStorage.getItem("ip");$.ajax({async:!0,type:"POST",dataType:"json",url:"/api",data:{info:"vsession",city:t,country:e,screen_size:o,page:n,battery:a,ip:c},success:function(){}}).done(function(t){})}function geoData(){navigator.geolocation?navigator.geolocation.getCurrentPosition(success,error,{enableHighAccuracy:!0,timeout:5e3,maximumAge:0}):x.innerHTML="Geolocation is not supported by this browser."}function success(t){var e=t.coords;localStorage.setItem("latitude",e.latitude),localStorage.setItem("longitude",e.longitude)}function error(t){console.warn(`ERROR(${t.code}): ${t.message}`)}