import React from 'react';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Switch from '@material-ui/core/Switch';

import IntervalActionsMenu from './IntervalActionsMenu';
import { formatTime, weekdaysShort } from '../util';

class IntervalListItem extends React.Component {
  render() {
    const { onToggle, onEdit, onDelete, interval } = this.props;

    return (
      <ListItem>
        <ListItemText
          primary={this.renderDayTimes()}
          secondary={this.renderWeekdays()}
        />
        <ListItemSecondaryAction>
          <Switch
            color="primary"
            checked={interval.enabled}
            onChange={onToggle}
          />
          <IntervalActionsMenu onEdit={onEdit} onDelete={onDelete} />
        </ListItemSecondaryAction>
      </ListItem>
    );
  }

  renderDayTimes() {
    const { from, to } = this.props.interval;

    return `${formatTime(from)} - ${formatTime(to)}`;
  }

  renderWeekdays() {
    const { weekdays } = this.props.interval;

    return weekdays.map(i => weekdaysShort[i]).join(', ');
  }
}

export default IntervalListItem;
