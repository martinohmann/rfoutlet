import React from 'react';
import List from '@material-ui/core/List';

import GroupListItem from './GroupListItem';
import OutletListItem from './OutletListItem';

class Group extends React.Component {
  state = {
    outlets: [],
  }

  componentDidMount() {
    this.updateState(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.updateState(nextProps);
  }

  updateState(props) {
    const { outlets } = props;

    this.setState({ outlets });
  }

  handleAction = action => event => {
    const { id } = this.props;

    this.props.dispatchMessage({ type: 'group', data: { id, action } });
  }

  render() {
    const { name } = this.props;
    const { outlets } = this.state;

    return (
      <List component="nav">
        <GroupListItem
          name={name}
          onActionOn={this.handleAction('on')}
          onActionOff={this.handleAction('off')}
          onActionToggle={this.handleAction('toggle')}
        />
        {outlets.map(outlet =>
          <OutletListItem key={outlet.id} {...this.props} {...outlet} />
        )}
      </List>
    );
  }
}

export default Group;
