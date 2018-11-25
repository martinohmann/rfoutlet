import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import List from '@material-ui/core/List';
import Divider from '@material-ui/core/Divider';

import DialogAppBar from './DialogAppBar';
import WeekdayListItem from './WeekdayListItem';
import { weekdaysLong } from '../util';

const styles = theme => ({
  container: {
    marginTop: 64,
  },
});

class WeekdaysDialog extends React.Component {
  state = {
    open: false,
    selected: [],
  }

  componentDidMount() {
    const { selected } = this.props;

    this.setState({ selected });
  }

  componentWillReceiveProps(nextProps) {
    const { open, selected } = nextProps;

    this.setState({ open, selected });
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
    const { classes, onClose } = this.props;
    const { open, selected } = this.state;

    return (
      <Dialog fullScreen open={open} onClose={onClose}>
        <DialogAppBar
          title="Select Weekdays"
          onClose={onClose}
          onDone={onClose}
          doneButtonDisabled={selected.length === 0}
          doneButtonText="Done"
        />
        <List component="nav" className={classes.container}>
          {weekdaysLong.map((weekday, key) => (
            <div key={key}>
              <WeekdayListItem
                weekday={weekday}
                selected={selected.indexOf(key) > -1}
                onToggle={this.handleWeekdayToggle(key)}
              />
              <Divider />
            </div>
          ))}
        </List>
      </Dialog>
    );
  }
}

WeekdaysDialog.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(WeekdaysDialog);
