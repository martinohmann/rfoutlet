import React from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from './List';
import ListItemText from '@material-ui/core/ListItemText';

import { formatDayTime, formatWeekdays } from '../util';

export default function IntervalOptionsList(props) {
  const {
    weekdays,
    fromDayTime,
    toDayTime,
    onWeekdaysClick,
    onFromDayTimeClick,
    onToDayTimeClick
  } = props;

  return (
    <List>
      <ListItem onClick={onWeekdaysClick}>
        <ListItemText primary="Weekdays" secondary={formatWeekdays(weekdays)} />
      </ListItem>
      <ListItem onClick={onFromDayTimeClick}>
        <ListItemText primary="From" secondary={formatDayTime(fromDayTime)} />
      </ListItem>
      <ListItem onClick={onToDayTimeClick}>
        <ListItemText primary="To" secondary={formatDayTime(toDayTime)} />
      </ListItem>
    </List>
  );
}

IntervalOptionsList.propTypes = {
    weekdays: PropTypes.array,
    fromDayTime: PropTypes.object,
    toDayTime: PropTypes.object,
    onWeekdaysClick: PropTypes.func.isRequired,
    onFromDayTimeClick: PropTypes.func.isRequired,
    onToDayTimeClick: PropTypes.func.isRequired,
};
