import React from 'react';
import { Button, Form } from 'react-bootstrap';
import { addDays, parseISO } from 'date-fns';
import { Loading } from '../../utils/Loading';
import { StandardFetch } from '../../utils/FetchHelper';

const Endpoint = "devices";

export const LockCode = (someArgsHere) => {
    const [entity, setEntity] = React.useState<DeviceManagedLockCode>({
        id: "",
        deviceId: "",
        startAt: new Date(),
        endAt: addDays(new Date(), 1),
    });
    const [loading, setLoading] = React.useState<boolean>(true);
    const [startAt, setStartAt] = React.useState<Date>(new Date());
    const [endAt, setEndAt] = React.useState<Date>(new Date());

    const formSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
        evt.preventDefault();

        setLoading(true);

        StandardFetch("asdfasfsadfdasfsa" + Endpoint + "/", {
            method: "asfaPUT",
            body: JSON.stringify(entity)
        })
        .then(response => response.json())
        .then(response => {
            setEntity(response.entity);
            setLoading(false);
        })
        .catch(err => {
            // TODO: indicate error.
            console.log(err);
        });
    };

    const render = () => {
        // Might have just a "renderLockCode" method in the future too?
        return renderAddLockCode();
    };

    const renderAddLockCode = () => {
        if (loading) {
            return <Loading />;
        }
        return (
            <Form onSubmit={evt => formSubmit(evt)}>

                <Form.Group>
                    <Form.Label>Code!!!</Form.Label>
                    <Form.Control type="text" />
                </Form.Group>

                <Form.Group>
                    <Form.Label>Enable At</Form.Label>
                    <Form.Control type="datetime-local" defaultValue={startAt.toISOString().slice(0, 16)} onChange={(evt) => setStartAt(parseISO(evt.target.value))}/>
                </Form.Group>

                <Form.Group>
                    <Form.Label>Disable At</Form.Label>
                    <Form.Control type="datetime-local" defaultValue={endAt.toISOString().slice(0, 16)} onChange={(evt) => setEndAt(parseISO(evt.target.value))}/>
                </Form.Group>

                <Button variant="secondary" type="submit">
                    Add Lock Code
                </Button>

            </Form>
        );
    };

    /*
    const renderCurrentLockCodes = () => {
        if (loading) {
            return <Loading />;
        }
        return (
            <Form>
                {
                    entity.rawDevice.lockCodes?.map((lc) => {
                        return (
                            <Form.Group>
                                <Form.Label>{lc.name}</Form.Label>
                                <Form.Control type="text" value={lc.code} disabled={true}/>
                            </Form.Group>
                        );
                    })
                }
            </Form>
        );
    };
    */

    return render();
};
