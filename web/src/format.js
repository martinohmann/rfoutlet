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

export function formatDayTime(dayTime, trans = ((k) => k)) {
  if (null === dayTime) {
    return trans('unset');
  }

  return dayTime.toFormat('HH:mm');
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
