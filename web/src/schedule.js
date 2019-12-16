import { DateTime } from 'luxon';

export const weekdaysLong = [
  'sunday',
  'monday',
  'tuesday',
  'wednesday',
  'thursday',
  'friday',
  'saturday',
];

export const weekdaysShort = [
  'sun',
  'mon',
  'tue',
  'wed',
  'thu',
  'fri',
  'sat',
];

export function formatTime(date) {
  return date.toFormat('HH:mm');
}

export function formatDayTime(dayTime, trans = ((k) => k)) {
  if (null === dayTime) {
    return trans('unset');
  }

  return formatTime(dayTime);
}

export function formatDayTimeInterval(interval, trans = ((k) => k)) {
  if (null === interval) {
    return trans('unset');
  }

  return `${formatDayTime(interval.from)} - ${formatDayTime(interval.to)}`;
}

export function formatWeekdays(weekdays, trans = ((k) => k)) {
  if (weekdays.length === 0) {
    return trans('unset');
  }

  return weekdays.map(i => trans(weekdaysShort[i])).join(', ');
}

export function formatSchedule(schedule, trans = ((k) => k)) {
  const intervals = schedule.filter(interval => interval.enabled);

  if (intervals.length === 0) {
    return '';
  }

  return trans('intervals-scheduled', { count: intervals.length });
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

function intervalToApp(interval) {
  const { from, to } = intervalToDateTimes(interval);

  return { ...interval, from, to };
}

export function intervalToApi(interval) {
  const { from, to } = dateTimesToInterval(interval);

  return { ...interval, from, to };
}

export function scheduleToApp(schedule) {
  const intervals = schedule || [];

  return intervals.map(interval => intervalToApp(interval));
}
