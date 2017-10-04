import React from 'react';

import OutletGroup from './OutletGroup';
import makeApiRequest from './api.js'

class OutletGroupContainer extends React.Component {
  constructor(props, context) {
    super(props, context)

    this.state = {
      outletGroups: props.outletGroups,
    };
  }

  componentDidMount() {
    makeApiRequest('/status', {}, result => {
      this.setState({ outletGroups: result });
    });
  }

  render() {
    return (
      <div style={this.props.style}>
        {this.state.outletGroups.map((attributes, key) =>
          <OutletGroup
            key={key}
            groupId={key}
            attributes={attributes}
          />
        )}
      </div>
    );
  }
}

export default OutletGroupContainer;
