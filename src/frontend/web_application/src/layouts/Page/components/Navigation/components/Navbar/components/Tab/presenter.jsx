import React, { Component } from 'react';
import PropTypes from 'prop-types';
import NavbarItem from '../NavbarItem';
import Icon from '../../../../../../../../components/Icon';
import ItemButton from '../ItemButton';
import ItemLink from '../ItemLink';

class Tab extends Component {
  static propTypes = {
    className: PropTypes.string,
    tab: PropTypes.shape({}).isRequired,
    onRemove: PropTypes.func.isRequired,
    isActive: PropTypes.bool.isRequired,
    last: PropTypes.bool.isRequired,
  };

  static defaultProps = {
    className: undefined,
  };

  constructor(props) {
    super(props);
    this.handleRemove = this.handleRemove.bind(this);
  }

  handleRemove() {
    this.props.onRemove(this.props.tab);
  }

  render() {
    const { className, tab, isActive, last } = this.props;

    return (
      <NavbarItem
        className={className}
        active={isActive}
        contentChildren={(
          <ItemLink to={tab.pathname}>
            <Icon type={tab.icon || 'dot-circle-o'} />
            {' '}
            {tab.label}
          </ItemLink>
        )}
        actionChildren={<ItemButton onClick={this.handleRemove}><Icon type="remove" /></ItemButton>}
        last={last}
      />
    );
  }
}

export default Tab;