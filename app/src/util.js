import config from './config';

export function apiRequest(method, requestUri, data = {}) {
  const url = config.api.baseUri + requestUri;

  const options = {
    method: method,
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
  };

  if ('POST' === method || 'PUT' === method) {
    options.body = JSON.stringify(data);
  }

  return fetch(url, options)
    .then(response => response.json());
}

export function formatTime(date) {
  return date.toFormat('HH:mm');
}
