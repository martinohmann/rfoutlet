import { scheduleToApp } from './schedule';

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
