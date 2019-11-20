import React from 'react';
import PropTypes from 'prop-types';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Checkbox from '@material-ui/core/Checkbox';

export default function WeekdayListItem(props) {
  const { onToggle, weekday, selected } = props;

  return (
    <ListItem onClick={onToggle}>
      <ListItemText primary={weekday} />
      <ListItemSecondaryAction>
        <Checkbox
          color="primary"
          onChange={onToggle}
          checked={selected}
        />
      </ListItemSecondaryAction>
    </ListItem>
  );
}

WeekdayListItem.propTypes = {
  onToggle: PropTypes.func.isRequired,
  selected: PropTypes.bool.isRequired,
  weekday: PropTypes.string.isRequired,
};
