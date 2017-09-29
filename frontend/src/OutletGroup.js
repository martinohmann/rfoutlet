import React from 'react';
import { List } from 'material-ui/List';
import FlatButton from 'material-ui/FlatButton';
import { Toolbar, ToolbarGroup, ToolbarTitle } from 'material-ui/Toolbar';

import Outlet from './Outlet';

class OutletGroup extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.handleOnButtonClick = this.handleOnButtonClick.bind(this);
    this.handleOffButtonClick = this.handleOffButtonClick.bind(this);
    this.handleToggleButtonClick = this.handleToggleButtonClick.bind(this);
    this.registerOutlet = this.registerOutlet.bind(this);

    this.styles = {
      toolbar: {
        paddingLeft: 16,
        paddingRight: 6,
      },
      title: {
        fontSize: 14,
      },
      button: {
        margin: 0,
        minWidth: 0,
      },
      label: {
        padding: 10,
        fontSize: 12,
      },
    };

    this.state = {
      outlets: [],
    }
  }

  registerOutlet(outlet) {
    var outlets = this.state.outlets;
    outlets.push(outlet);
    this.setState({ outlets: outlets });
  }

  updateOutletStates(callback) {
    this.state.outlets.map(outlet => {
      return outlet.setState({ isEnabled: callback(outlet) });
    });
  }

  handleOnButtonClick() {
    this.updateOutletStates(outlet => true);
  }

  handleOffButtonClick() {
    this.updateOutletStates(outlet => false);
  }

  handleToggleButtonClick() {
    this.updateOutletStates(outlet => !outlet.state.isEnabled);
  }

  render() {
    return (
      <List>
        <Toolbar style={this.styles.toolbar}>
          <ToolbarGroup>
            <ToolbarTitle text={this.props.identifier} style={this.styles.title} />
          </ToolbarGroup>
          <ToolbarGroup>
            <FlatButton
              label="On"
              primary={true}
              style={this.styles.button}
              labelStyle={this.styles.label}
              onClick={this.handleOnButtonClick}
            />
            <FlatButton
              label="Off"
              secondary={true}
              style={this.styles.button}
              labelStyle={this.styles.label}
              onClick={this.handleOffButtonClick}
            />
            <FlatButton
              label="Toggle"
              style={this.styles.button}
              labelStyle={this.styles.label}
              onClick={this.handleToggleButtonClick}
            />
          </ToolbarGroup>
        </Toolbar>
        {this.props.outlets.map((outlet, key) =>
          <Outlet
            key={key}
            identifier={outlet.identifier}
            state={outlet.state}
            registerOutlet={this.registerOutlet}
          />
        )}
      </List>
    );
  }
}

export default OutletGroup;

