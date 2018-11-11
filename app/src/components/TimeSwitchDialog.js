import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Button from '@material-ui/core/Button';
import { TimePicker } from 'material-ui-pickers';

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  formControl: {
    margin: theme.spacing.unit,
  },
});

class TimeSwitchDialog extends React.Component {
  state = {
    from: null,
    to: null,
    open: false,
  }

  componentWillReceiveProps(nextProps) {
    const { open, from, to } = nextProps;

    this.setState({ open, from, to });
  }

  handleChange = name => date => {
    this.setState({ [name]: date });
  }

  handleClear = () => {
    this.setState({ from: null, to: null });
  }

  handleClose = () => {
    this.props.onClose();
  }

  handleApply = () => {
    const { from, to } = this.state;
    const enabled = (from !== null && to !== null);

    this.props.onApply({ from, to, enabled });
  }

  render() {
    const { classes, identifier } = this.props;
    const { open, from, to } = this.state;

    return (
      <Dialog open={open} onClose={this.handleClose}>
        <DialogTitle>Timer for {identifier}</DialogTitle>
        <DialogContent className={classes.container}>
          <TimePicker
            className={classes.formControl}
            clearable
            ampm={false}
            label="From"
            value={from}
            onChange={this.handleChange('from')}
          />
          <TimePicker
            className={classes.formControl}
            clearable
            ampm={false}
            label="To"
            value={to}
            onChange={this.handleChange('to')}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={this.handleClear} color="primary">
            Clear
          </Button>
          <Button onClick={this.handleClose} color="primary">
            Cancel
          </Button>
          <Button onClick={this.handleApply} color="secondary">
            Apply
          </Button>
        </DialogActions>
      </Dialog>
    );
  }
}

TimeSwitchDialog.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(TimeSwitchDialog);
