import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Icon from '../../../Icon';
import './style.scss';

class AddressDetails extends Component {
  static propTypes = {
    address: PropTypes.shape({}).isRequired,
    __: PropTypes.func.isRequired,
  };

  constructor(props) {
    super(props);
    this.initTranslations();
  }

  initTranslations() {
    const { __ } = this.props;
    this.addressTypesTranslations = {
      work: __('contact.address_type.work'),
      home: __('contact.address_type.home'),
      other: __('contact.address_type.other'),
    };
  }

  render() {
    const { address } = this.props;

    return (
      <span className="m-address-details">
        <Icon type="map-marker" rightSpaced />
        <address className="m-address-details__postal-address">
          {address.street}{address.street && ', '}
          {address.postal_code}{address.postal_code && ' '}
          {address.city}{address.city && ' '}
          {address.country}{address.country && ' '}
          {address.region}
        </address>
        {' '}
        {(address.label || address.type) &&
          <small>
            <em>({address.label}{address.label && ' '}{this.addressTypesTranslations[address.type]})</em>
          </small>
        }
      </span>
    );
  }
}

export default AddressDetails;
