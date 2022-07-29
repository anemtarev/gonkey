package parser

import (
    "errors"
    "fmt"
    "io/ioutil"

    "gopkg.in/yaml.v3"
)

type redisYamlParser struct {
    fileParser *fileParser
}

func (p *redisYamlParser) Copy(fileParser *fileParser) FixtureFileParser {
    cp := &(*p)
    cp.fileParser = fileParser
    return cp
}

func (p *redisYamlParser) buildTemplate(ctx *context, f Fixture) error {
    for _, tplData := range f.Templates.Keys {
        refName := tplData.Name
        if refName == "" {
            return errors.New("template $name is required")
        }
        if _, ok := ctx.keyRefs[refName]; ok {
            return fmt.Errorf("unable to load template %s: duplicating ref name", refName)
        }
        if tplData.Extend != "" {
            baseRecord, err := p.resolveKeyReference(ctx.keyRefs, tplData.Extend)
            if err != nil {
                return err
            }
            for k, v := range tplData.Values {
                baseRecord.Values[k] = v
            }
            tplData.Values = baseRecord.Values
        }

        keyRef := Keys{
            Values: make(map[string]*KeyValue, len(tplData.Values)),
        }
        for k, v := range tplData.Values {
            keyRef.Values[k] = v
        }
        ctx.keyRefs[refName] = keyRef
    }

    for _, tplData := range f.Templates.Sets {
        refName := tplData.Name
        if refName == "" {
            return errors.New("template $name is required")
        }
        if _, ok := ctx.setRefs[refName]; ok {
            return fmt.Errorf("unable to load template %s: duplicating ref name", refName)
        }
        if tplData.Extend != "" {
            baseRecord, err := p.resolveSetReference(ctx.setRefs, tplData.Extend)
            if err != nil {
                return err
            }
            for k, v := range tplData.Values {
                baseRecord.Values[k] = v
            }
            tplData.Values = baseRecord.Values
        }

        setRef := SetRecordValue{
            Values: make(map[string]*SetValue),
        }
        for k, v := range tplData.Values {
            var valueCopy *SetValue
            if v != nil {
                valueCopy = &(*v)
            }
            setRef.Values[k] = valueCopy
        }
        ctx.setRefs[refName] = setRef
    }

    for _, tplData := range f.Templates.Maps {
        refName := tplData.Name
        if refName == "" {
            return errors.New("template $name is required")
        }
        if _, ok := ctx.mapRefs[refName]; ok {
            return fmt.Errorf("unable to load template %s: duplicating ref name", refName)
        }
        if tplData.Extend != "" {
            baseRecord, err := p.resolveMapReference(ctx.mapRefs, tplData.Extend)
            if err != nil {
                return err
            }
            for k, v := range tplData.Values {
                baseRecord.Values[k] = v
            }
            tplData.Values = baseRecord.Values
        }
        mapRef := MapRecordValue{
            Values: make(map[string]string, len(tplData.Values)),
        }
        for k, v := range tplData.Values {
            mapRef.Values[k] = v
        }
        ctx.mapRefs[refName] = mapRef
    }

    return nil
}

func (p *redisYamlParser) resolveKeyReference(refs map[string]Keys, refName string) (*Keys, error) {
    refTemplate, ok := refs[refName]
    if !ok {
        return nil, fmt.Errorf("ref not found: %s", refName)
    }
    keysCopy := &Keys{
        Values: make(map[string]*KeyValue),
    }
    for k, v := range refTemplate.Values {
        var valueCopy *KeyValue
        if v != nil {
            valueCopy = &(*v)
        }
        keysCopy.Values[k] = valueCopy
    }
    return keysCopy, nil
}

func (p *redisYamlParser) resolveMapReference(refs map[string]MapRecordValue, refName string) (*MapRecordValue, error) {
    refTemplate, ok := refs[refName]
    if !ok {
        return nil, fmt.Errorf("ref not found: %s", refName)
    }
    mapCopy := &MapRecordValue{
        Values: make(map[string]string),
    }
    for k, v := range refTemplate.Values {
        mapCopy.Values[k] = v
    }
    return mapCopy, nil
}

func (p *redisYamlParser) resolveSetReference(refs map[string]SetRecordValue, refName string) (*SetRecordValue, error) {
    refTemplate, ok := refs[refName]
    if !ok {
        return nil, fmt.Errorf("ref not found: %s", refName)
    }
    setCopy := &SetRecordValue{
        Values: make(map[string]*SetValue),
    }
    for k, v := range refTemplate.Values {
        var setValueCopy *SetValue
        if v != nil {
            setValueCopy = &(*v)
        }
        setCopy.Values[k] = setValueCopy
    }
    return setCopy, nil
}

func (p *redisYamlParser) buildKeys(ctx *context, data *Keys) error {
    if data == nil {
        return nil
    }
    if data.Extend != "" {
        baseRecord, err := p.resolveKeyReference(ctx.keyRefs, data.Extend)
        if err != nil {
            return err
        }
        for k, v := range data.Values {
            var keyValueCopy *KeyValue
            if v != nil {
                keyValueCopy = &(*v)
            }
            baseRecord.Values[k] = keyValueCopy
        }
        data.Values = baseRecord.Values
    }
    return nil
}

func (p *redisYamlParser) buildMaps(ctx *context, data *Maps) error {
    if data == nil {
        return nil
    }
    for _, v := range data.Values {
        if v.Extend != "" {
            baseRecord, err := p.resolveMapReference(ctx.mapRefs, v.Extend)
            if err != nil {
                return err
            }
            for k, v := range v.Values {
                baseRecord.Values[k] = v
            }
            v.Values = baseRecord.Values
        }
        if v.Name != "" {
            mapRef := MapRecordValue{
                Values: make(map[string]string, len(v.Values)),
            }
            for k, v := range v.Values {
                mapRef.Values[k] = v
            }
            ctx.mapRefs[v.Name] = mapRef
        }
    }
    return nil
}

func (p *redisYamlParser) buildSets(ctx *context, data *Sets) error {
    if data == nil {
        return nil
    }
    for _, v := range data.Values {
        if v.Extend != "" {
            baseRecord, err := p.resolveSetReference(ctx.setRefs, v.Extend)
            if err != nil {
                return err
            }
            for k, v := range v.Values {
                baseRecord.Values[k] = v
            }
            v.Values = baseRecord.Values
        }
        if v.Name != "" {
            setRef := SetRecordValue{
                Values: make(map[string]*SetValue),
            }
            for k, v  := range v.Values {
                var setValueCopy *SetValue
                if v != nil {
                    setValueCopy = &(*v)
                }
                setRef.Values[k] = setValueCopy
            }
            ctx.setRefs[v.Name] = setRef
        }
    }
    return nil
}

func (p *redisYamlParser) Parse(ctx *context, filename string) (*Fixture, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var fixture Fixture
    if err := yaml.Unmarshal(data, &fixture); err != nil {
        return nil, err
    }

    for _, parentFixture := range fixture.Inherits {
        _, err := p.fileParser.ParseFiles(ctx, []string{parentFixture})
        if err != nil {
            return nil, err
        }
    }

    if err = p.buildTemplate(ctx, fixture); err != nil {
        return nil, err
    }

    for _, databaseData := range fixture.Databases {
        if err := p.buildKeys(ctx, databaseData.Keys); err != nil {
            return nil, err
        }
        if err := p.buildMaps(ctx, databaseData.Maps); err != nil {
            return nil, err
        }
        if err := p.buildSets(ctx, databaseData.Sets); err != nil {
            return nil, err
        }
    }

    return &fixture, nil
}
