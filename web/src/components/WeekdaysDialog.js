import React from 'react';
import PropTypes from 'prop-types';
import { List } from './List';

import ConfigurationDialog from './ConfigurationDialog';
import WeekdayListItem from './WeekdayListItem';
import { weekdaysLong } from '../util';

class WeekdaysDialog extends React.Component {
  state = {
    selected: [],
  }

  componentDidMount() {
    const { selected } = this.props;

    this.setState({ selected });
  }

  handleWeekdayToggle = key => () => {
    this.setState(state => {
      const { selected } = state;
      const index = selected.indexOf(key);

      if (index > -1) {
        selected.splice(index, 1);
      } else {
        selected.push(key);
      }

      selected.sort();

      return { selected };
    });
  }

  render() {
    const { open, onClose } = this.props;
    const { selected } = this.state;

    return (
      <ConfigurationDialog
        title="Select Weekdays"
        open={open}
        onClose={onClose}
        onDone={onClose}
        doneButtonDisabled={selected.length === 0}
        doneButtonText="Done"
      >
        <List>
          {weekdaysLong.map((weekday, key) => (
            <WeekdayListItem
              key={key}
              weekday={weekday}
              selected={selected.indexOf(key) > -1}
              onToggle={this.handleWeekdayToggle(key)}
            />
          ))}
        </List>
      </ConfigurationDialog>
    );
  }
}

WeekdaysDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
};

export default WeekdaysDialog;
