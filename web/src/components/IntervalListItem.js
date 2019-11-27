import React from 'react';
import PropTypes from 'prop-types';
import { ListItem } from './List';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Switch from '@material-ui/core/Switch';

import IntervalActionsMenu from './IntervalActionsMenu';
import { formatDayTimeInterval, formatWeekdays } from '../schedule';

export default function IntervalListItem(props) {
  const { interval, onDelete, onEdit, onToggle } = props;

  return (
    <ListItem onClick={onEdit}>
      <ListItemText
        primary={formatDayTimeInterval(interval)}
        secondary={formatWeekdays(interval.weekdays)}
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

IntervalListItem.propTypes = {
    interval: PropTypes.object,
    onDelete: PropTypes.func.isRequired,
    onEdit: PropTypes.func.isRequired,
    onToggle: PropTypes.func.isRequired,
};
