import { DateTime } from 'luxon';

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

function scheduleToApp(schedule) {
  const intervals = schedule || [];

  return intervals.map(interval => intervalToApp(interval));
}

export function convertToApp(groups) {
  return groups.map(group => {
    const outlets = group.outlets || [];

    group.outlets = outlets.map(outlet => {
      outlet.schedule = scheduleToApp(outlet.schedule);

      return outlet;
    });

    return group;
  });
}

export function intervalToApi(interval) {
  const { from, to } = dateTimesToInterval(interval);

  return { ...interval, from, to };
}
