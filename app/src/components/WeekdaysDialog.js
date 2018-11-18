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

  handleDone = () => {
    const { selected } = this.state;

    this.props.onDone(selected);
  }

  handleToggle = value => () => {
    this.setState(state => {
      const { selected } = state;
      const index = selected.indexOf(value);

      if (index > -1) {
        selected.splice(index, 1);
      } else {
        selected.push(value);
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
          {weekdaysLong.map((name, value) => (
            <div key={value}>
              <WeekdayListItem
                weekday={name}
                selected={selected.indexOf(value) > -1}
                onToggle={this.handleToggle(value)}
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
