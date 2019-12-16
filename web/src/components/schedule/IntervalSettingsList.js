import React from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from '../List';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import DateRangeIcon from '@material-ui/icons/DateRange';
import TimerIcon from '@material-ui/icons/Timer';
import TimerOffIcon from '@material-ui/icons/TimerOff';
import { useTranslation } from 'react-i18next';
import { formatDayTime, formatWeekdays } from '../../schedule';

export default function IntervalSettingsList(props) {
  const {
    weekdays,
    fromDayTime,
    toDayTime,
    onWeekdaysClick,
    onFromDayTimeClick,
    onToDayTimeClick
  } = props;

  const { t } = useTranslation();

  return (
    <List>
      <ListItem onClick={onWeekdaysClick}>
        <ListItemIcon>
          <DateRangeIcon />
        </ListItemIcon>
        <ListItemText primary={t('weekdays')} secondary={formatWeekdays(weekdays, t)} />
      </ListItem>
      <ListItem onClick={onFromDayTimeClick}>
        <ListItemIcon>
          <TimerIcon />
        </ListItemIcon>
        <ListItemText primary={t('start-time')} secondary={formatDayTime(fromDayTime, t)} />
      </ListItem>
      <ListItem onClick={onToDayTimeClick}>
        <ListItemIcon>
          <TimerOffIcon />
        </ListItemIcon>
        <ListItemText primary={t('end-time')} secondary={formatDayTime(toDayTime, t)} />
      </ListItem>
    </List>
  );
}

IntervalSettingsList.propTypes = {
    weekdays: PropTypes.array,
    fromDayTime: PropTypes.object,
    toDayTime: PropTypes.object,
    onWeekdaysClick: PropTypes.func.isRequired,
    onFromDayTimeClick: PropTypes.func.isRequired,
    onToDayTimeClick: PropTypes.func.isRequired,
};
