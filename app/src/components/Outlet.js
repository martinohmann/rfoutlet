import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Switch from '@material-ui/core/Switch';

import { apiRequest, outletEnabled } from '../util'

const styles = {};

class Outlet extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.handleToggle = this.handleToggle.bind(this)

    this.state = {
      isEnabled: outletEnabled(props.outlet),
    };

    this.props.registerOutlet(props.outlet.id, this);
  }

  componentWillReceiveProps(nextProps) {
    const outlet = nextProps.outlet;

    this.setState({ isEnabled: outletEnabled(outlet) });
  }

  handleToggle(event, isEnabled) {
    const data = {
      action: 'toggle',
      id: this.props.outlet.id
    };

    apiRequest('POST', '/outlet', data)
      .then(outlet => this.setState({ isEnabled: outletEnabled(outlet) }))
      .catch(err => console.error(err));
  }

  render() {
    return (
      <ListItem button onClick={this.handleToggle}>
        <ListItemText primary={this.props.outlet.name} />
        <ListItemSecondaryAction>
          <Switch onChange={this.handleToggle} checked={this.state.isEnabled}
          />
        </ListItemSecondaryAction>
      </ListItem>
    );
  }
}

Outlet.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Outlet);
