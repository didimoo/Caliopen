import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withTranslator } from '@gandi/react-translate';
import { Field } from 'redux-form';
import ContactTitleField from '../ContactTitleField';
import renderReduxField from '../../../../services/renderReduxField';
import Button from '../../../../components/Button';
import { TextFieldGroup as TextFieldGroupBase } from '../../../../components/form';
import './style.scss';

const TextFieldGroup = renderReduxField(TextFieldGroupBase);

class ContactProfileForm extends Component {
  static propTypes = {
    __: PropTypes.func.isRequired,
    form: PropTypes.string.isRequired,
    isNew: PropTypes.bool,
  };

  static defaultProps = {
    isNew: false,
  };

  state = {
    isExpanded: false,
  };

  componentWillMount() {
    if (this.props.isNew) {
      this.setState({
        isExpanded: true,
      });
    }
  }

  toggleExpandForm = () => {
    this.setState(prevState => ({
      ...prevState,
      isExpanded: !prevState.isExpanded,
    }));
  };

  render() {
    const { __, form } = this.props;

    return (
      <div className="m-contact-profile-form">
        <div className="m-contact-profile-form__header">
          <ContactTitleField className="m-contact-profile-form__title" form={form} />
          {this.state.isExpanded ?
            <Button
              icon="caret-up"
              display="inline"
              onClick={this.toggleExpandForm}
              className="m-contact-profile-form__expand-button"
            >
              <span className="show-for-sr">
                {__('contact_profile.action.edit_contact')}
              </span>
            </Button>
          :
            <Button
              icon="caret-down"
              display="inline"
              onClick={this.toggleExpandForm}
              className="m-contact-profile-form__expand-button"
            >
              <span className="show-for-sr">
                {__('contact_profile.action.edit_contact')}
              </span>
            </Button>
          }
        </div>

        {this.state.isExpanded &&
          <div className="m-contact-profile-form__expanded-form">
            <Field
              component={TextFieldGroup}
              className="m-contact-profile-form__input"
              label={__('contact_profile.form.name-prefix.label')}
              placeholder={__('contact_profile.form.name-prefix.label')}
              name="name_prefix"
              showLabelforSr
            />
            <Field
              component={TextFieldGroup}
              className="m-contact-profile-form__input"
              label={__('contact_profile.form.firstname.label')}
              placeholder={__('contact_profile.form.firstname.label')}
              name="given_name"
              showLabelforSr
            />
            <Field
              component={TextFieldGroup}
              className="m-contact-profile-form__input"
              label={__('contact_profile.form.lastname.label')}
              placeholder={__('contact_profile.form.lastname.label')}
              name="family_name"
              showLabelforSr
            />
            <Field
              component={TextFieldGroup}
              className="m-contact-profile-form__input"
              label={__('contact_profile.form.name-suffix.label')}
              placeholder={__('contact_profile.form.name-suffix.label')}
              name="name_suffix"
              showLabelforSr
            />
          </div>
        }
      </div>
    );
  }
}


export default withTranslator()(ContactProfileForm);
