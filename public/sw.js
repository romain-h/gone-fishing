var CACHE = 'gone-fishing-wpa-1-3';
var filesToCache = [
  '/',
  '/assets/main.css',
  '/assets/favicon-32x32.png',
  '/assets/favicon-16x16.png',
];

self.addEventListener('install', function(e) {
  console.log('[ServiceWorker] Install');
  e.waitUntil(
    caches.open(CACHE).then(function(cache) {
      console.log('[ServiceWorker] Caching app shell');
      cache.add('//stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css');
      cache.add('//code.jquery.com/jquery-3.3.1.slim.min.js');
      cache.add('//stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js');
      return cache.addAll(filesToCache);
    })
  );
});

self.addEventListener('activate', function(e) {
  console.log('[ServiceWorker] Activate');
  e.waitUntil(
    caches.keys().then(function(keyList) {
      return Promise.all(keyList.map(function(key) {
        if (key !== CACHE) {
          console.log('[ServiceWorker] Removing old cache', key);
          return caches.delete(key);
        }
      }));
    })
  );
  return self.clients.claim();
});

self.addEventListener('fetch', function(e) {
  console.log('[ServiceWorker] Fetch', e.request.url);
  e.respondWith(fromCache(e.request));
  e.waitUntil(update(e.request).then(refresh));
});

self.addEventListener('foreignfetch', event => {
  event.respondWith(fetch(event.request).then(response => {
    return {
      response: response,
      origin: event.origin,
      headers: ['Content-Type']
    }
  }));
});

function fromCache(request) {
  return caches.open(CACHE).then(function (cache) {
    return cache.match(request);
  }).then(response => {
    return response || update(request);
  })
}

function update(request) {
  return caches.open(CACHE).then(function (cache) {
    return fetch(request).then(function (response) {
      return cache.put(request, response.clone()).then(function () {
        return response;
      });
    });
  });
}

function refresh(response) {
  return self.clients.matchAll().then(function (clients) {
    clients.forEach(function (client) {
      var message = {
        type: 'refresh',
        url: response.url,
        eTag: response.headers.get('ETag')
      };
      client.postMessage(JSON.stringify(message));
    });
  });
}
