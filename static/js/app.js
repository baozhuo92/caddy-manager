(function() {

  function showToast(message, type) {
    var container = document.getElementById('toastContainer');
    if (!container) {
      container = document.createElement('div');
      container.id = 'toastContainer';
      container.className = 'toast-container';
      document.body.appendChild(container);
    }
    var toast = document.createElement('div');
    toast.className = 'toast toast-' + (type || 'info');
    toast.textContent = message;
    container.appendChild(toast);
    setTimeout(function() {
      toast.remove();
    }, 3000);
  }

  function showLoadingMask() {
    var mask = document.getElementById('loadingMask');
    if (mask) mask.classList.add('active');
  }

  function hideLoadingMask() {
    var mask = document.getElementById('loadingMask');
    if (mask) mask.classList.remove('active');
  }

  window.showToast = showToast;
  window.showLoadingMask = showLoadingMask;
  window.hideLoadingMask = hideLoadingMask;

  var token = localStorage.getItem('auth_token');
  var currentPath = window.location.pathname;

  if (!token && currentPath !== '/login' && currentPath !== '/init') {
    window.location.href = '/login?redirect=' + encodeURIComponent(currentPath + window.location.search);
    return;
  }

  if (token && (currentPath === '/login' || currentPath === '/init')) {
    window.location.href = '/';
    return;
  }

  var originalFetch = window.fetch;
  window.fetch = function(url, options) {
    options = options || {};
    options.headers = options.headers || {};

    var t = localStorage.getItem('auth_token');
    if (t && !options.headers['Token']) {
      options.headers['Token'] = t;
    }

    return originalFetch(url, options).then(function(response) {
      if (response.status === 401) {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
        if (response.headers.get('Content-Type') && response.headers.get('Content-Type').indexOf('application/json') > -1) {
          window.location.href = '/login?redirect=' + encodeURIComponent(window.location.pathname + window.location.search);
        }
      }
      return response;
    });
  };

  if (token) {
    var username = localStorage.getItem('auth_user');
    if (!document.querySelector('.header-user span') && username) {
      var headerUser = document.querySelector('.header-user');
      if (headerUser) {
        var span = document.createElement('span');
        span.textContent = username;
        headerUser.insertBefore(span, headerUser.firstChild);
      }
    }
  }

})();
