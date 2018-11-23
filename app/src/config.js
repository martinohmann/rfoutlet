export default {
  project: {
    name: "rfoutlet",
    url: 'https://github.com/martinohmann/rfoutlet',
  },
  ws: {
    url: (() => {
      const l = window.location;

      if (!process.env.NODE_ENV || process.env.NODE_ENV === 'development') {
        return `ws://${l.hostname}:3333/ws`
      }

      return `ws://${l.host}/ws`;
    })(),
  }
}
