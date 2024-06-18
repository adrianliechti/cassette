let events = [];

rrweb.record({
  emit(event) {
    events.push(event);
  },
});

function save() {
  const body = JSON.stringify({ events });
  events = [];
  fetch('https://relevant-new-vervet.ngrok-free.app/events', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body,
  });
}

setInterval(save, 10 * 1000);
