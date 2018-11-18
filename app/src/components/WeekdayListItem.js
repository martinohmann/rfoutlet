import React from 'react';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Checkbox from '@material-ui/core/Checkbox';

class WeekdayListItem extends React.Component {
  render() {
    const { onToggle, weekday, selected } = this.props;

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
}

export default WeekdayListItem;
