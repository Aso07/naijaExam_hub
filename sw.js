const CACHE_NAME = 'naija-exam-v1';
const ASSETS = [
  '/',
  '/index.html',
  // We will cache the questions API response too
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(ASSETS))
  );
});

self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      // Return cached file if found, otherwise go to network
      return response || fetch(event.request);
    })
  );
});
