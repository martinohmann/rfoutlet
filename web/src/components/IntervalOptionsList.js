import React from 'react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Divider from '@material-ui/core/Divider';

import { formatTime, weekdaysShort } from '../util';

class IntervalOptionsList extends React.Component {
  render() {
    const {
      className,
      weekdays,
      fromDayTime,
      toDayTime,
      onWeekdaysClick,
      onFromDayTimeClick,
      onToDayTimeClick
    } = this.props;

    return (
      <List component="nav" className={className}>
        <ListItem onClick={onWeekdaysClick}>
          <ListItemText
            primary="Weekdays"
            secondary={this.renderWeekdays(weekdays)}
          />
        </ListItem>
        <Divider />
        <ListItem onClick={onFromDayTimeClick}>
          <ListItemText primary="From" secondary={this.renderDayTime(fromDayTime)} />
        </ListItem>
        <Divider />
        <ListItem onClick={onToDayTimeClick}>
          <ListItemText primary="To" secondary={this.renderDayTime(toDayTime)} />
        </ListItem>
        <Divider />
      </List>
    );
  }

  renderDayTime(dayTime) {
    if (null === dayTime) {
      return 'unset';
    }

    return formatTime(dayTime);
  }

  renderWeekdays(weekdays) {
    if (weekdays.length === 0) {
      return 'none';
    }

    return weekdays.map(i => weekdaysShort[i]).join(', ');
  }
}

export default IntervalOptionsList;
