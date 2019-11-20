import { DateTime } from 'luxon';

export function formatTime(date) {
  return date.toFormat('HH:mm');
}

export const weekdaysLong = [
  'Sunday',
  'Monday',
  'Tuesday',
  'Wednesday',
  'Thursday',
  'Friday',
  'Saturday',
];

export const weekdaysShort = [
  'Sun',
  'Mon',
  'Tue',
  'Wed',
  'Thu',
  'Fri',
  'Sat',
];

function intervalToDateTimes({ from, to }) {
  return {
    from: dayTimeToDateTime(from),
    to: dayTimeToDateTime(to),
  };
}

function dayTimeToDateTime(dayTime) {
  const { hour, minute } = dayTime;

  return DateTime.local().set({ hour, minute });
}

function dateTimesToInterval({ from, to }) {
  return {
    from: dateTimeToDayTime(from),
    to: dateTimeToDayTime(to),
  };
}

function dateTimeToDayTime(dateTime) {
  const { hour, minute } = dateTime;

  return { hour, minute };
}

export function scheduleToApp(schedule) {
  return (schedule || []).map(interval => {
    const { from, to } = intervalToDateTimes(interval);

    return { ...interval, from, to };
  });
}

export function intervalToApi(interval) {
  const { from, to } = dateTimesToInterval(interval);

  return { ...interval, from, to };
}
