let events = [];
let eventsURL = new URL('/events', document.currentScript.src).href;

rrweb.record({
  emit(event) {
    events.push(event);
  },
});

function save() {
  const body = JSON.stringify({ events });
  events = [];

  fetch(eventsURL, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });
}

setInterval(save, 10 * 1000);