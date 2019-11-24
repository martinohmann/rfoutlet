import { DateTime } from 'luxon';

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

export function formatTime(date) {
  return date.toFormat('HH:mm');
}

export function formatDayTime(dayTime) {
  if (null === dayTime) {
    return 'unset';
  }

  return formatTime(dayTime);
}

export function formatDayTimeInterval(interval) {
  if (null === interval) {
    return 'unset';
  }

  return `${formatDayTime(interval.from)} - ${formatDayTime(interval.to)}`;
}

export function formatWeekdays(weekdays) {
  if (weekdays.length === 0) {
    return 'none';
  }

  return weekdays.map(i => weekdaysShort[i]).join(', ');
}

export function formatSchedule(schedule) {
  const intervals = schedule.filter(interval => interval.enabled);

  if (intervals.length === 0) {
    return '';
  }

  if (intervals.length === 1) {
    return `1 interval scheduled`;
  }

  return `${intervals.length} intervals scheduled`;
}

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
