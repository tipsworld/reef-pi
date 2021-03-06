import React from 'react'
import Enzyme, {shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import EditTimer from './edit_timer'
import 'isomorphic-fetch'
import * as Alert from '../utils/alert'

Enzyme.configure({ adapter: new Adapter() })

describe('<EditTimer />', () => {
  var values = {type: 'equipment'}
  var equipment = [{id: '1', name: 'EQ'}]
  var fn = jest.fn()

  beforeEach(() => {
    jest.spyOn(Alert, 'showError')
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  it('<EditTimer />', () => {
    shallow(
      <EditTimer values={values}
        errors={{}}
        touched={{}}
        equipment={equipment}
        handleBlur={fn}
        handleChange={fn}
        submitForm={fn} />
    )
  })

  it('<EditTimer /> should submit', () => {
    const wrapper = shallow(
      <EditTimer values={values}
        equipment={equipment}
        handleBlur={fn}
        handleChange={fn}
        submitForm={fn}
        errors={{}}
        touched={{}}
        dirty
        isValid
        showChart={false} />
    )
    wrapper.find('form').simulate('submit', {preventDefault: () => {}})
    expect(Alert.showError).not.toHaveBeenCalled()
  })

  it('<EditTimer /> should show alert when invalid', () => {
    values.type = 'reminder'
    values.name = ''

    const wrapper = shallow(
      <EditTimer values={values}
        equipment={equipment}
        handleBlur={fn}
        handleChange={fn}
        submitForm={fn}
        showChart
        errors={{}}
        touched={{}}
        dirty
        isValid={false} />
    )
    wrapper.find('form').simulate('submit', {preventDefault: () => {}})
    expect(Alert.showError).toHaveBeenCalled()
  })
})
