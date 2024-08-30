(function() {
    const localStorage = window.localStorage;

    const propSeparator = ",,";

    function mangle(propKey) {
        return "$" + propKey + "$";
    }

    function unmangle(propKey) {
        return propKey.substring(1, propKey.length-1);
    }

    function objectPropPath(appKey, objectKey, mangledPropKey) {
        return appKey + "_" + objectKey + "__" + mangledPropKey;
    }

    function objectMetadataPath(appKey, objectKey) {
        return appKey + "_" + objectKey + "_proplist_";
    }

    function listObjectProps(appKey, objectKey) {
        let metaKey = objectMetadataPath(appKey, objectKey);
        let v = localStorage.getItem(metaKey);
        if (v === null) {
            return null;
        }
        if (v === "") {
            return [];
        }
        return v.split(propSeparator).map(unmangle);
    }

    function saveObjectProp(appKey, objectKey, propKey, data) {
        let mangledPropKey = mangle(propKey);
        let metaKey = objectMetadataPath(appKey, objectKey);
        let v = localStorage.getItem(metaKey);
        if (v === null || v === '') {
            // Creating a fresh meta file with a single (new) prop key.
            localStorage.setItem(metaKey, mangledPropKey);
        } else {
            if (!v.includes(mangledPropKey)) {
                v += propSeparator + mangledPropKey;
                localStorage.setItem(metaKey, v);
            }
        }
        localStorage.setItem(objectPropPath(appKey, objectKey, mangledPropKey), data);
    }

    function loadObjectProp(appKey, objectKey, propKey) {
        let mangledPropKey = mangle(propKey);
        return localStorage.getItem(objectPropPath(appKey, objectKey, mangledPropKey));
    }

    function objectPropExists(appKey, objectKey, propKey) {
        let mangledPropKey = mangle(propKey);
        return localStorage.getItem(objectPropPath(appKey, objectKey, mangledPropKey)) !== null;
    }

    function objectExists(appKey, objectKey) {
        return localStorage.getItem(objectMetadataPath(appKey, objectKey)) !== null;
    }

    function deleteObjectProp(appKey, objectKey, propKey) {
        let mangledPropKey = mangle(propKey);
        localStorage.removeItem(objectPropPath(appKey, objectKey, mangledPropKey));

        let metaKey = objectMetadataPath(appKey, objectKey);
        let v = localStorage.getItem(metaKey);
        if (v !== null) {
            let parts = v.split(propSeparator);
            let index = parts.indexOf(mangledPropKey);
            if (index > -1) {
                parts.splice(index, 1);
                localStorage.setItem(metaKey, parts.join(propSeparator));
            }
        }
    }

    function deleteObject(appKey, objectKey) {
        let metaKey = objectMetadataPath(appKey, objectKey);
        let v = localStorage.getItem(metaKey);
        if (v === null) {
            return;
        }

        let keys = v.split(propSeparator);
        for (let i = 0; i < keys.length; i++) {
            let mangledPropKey = keys[i];
            localStorage.removeItem(objectPropPath(appKey, objectKey, mangledPropKey));
        }
        localStorage.removeItem(metaKey);
    }

    return {
        "objectPropPath": objectPropPath,
        "listObjectProps": listObjectProps,
        "saveObjectProp": saveObjectProp,
        "loadObjectProp": loadObjectProp,
        "objectPropExists": objectPropExists,
        "objectExists": objectExists,
        "deleteObjectProp": deleteObjectProp,
        "deleteObject": deleteObject,
    };
}());
