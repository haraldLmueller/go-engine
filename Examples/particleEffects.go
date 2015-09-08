package examples

import (
    "math"
	"image/color"

    "github.com/Walesey/goEngine/vectorMath"
    "github.com/Walesey/goEngine/assets"
    "github.com/Walesey/goEngine/effects"
    "github.com/Walesey/goEngine/renderer"

    "github.com/codegangsta/cli"
)


//
func Particles( c *cli.Context ){
    fps := renderer.CreateFPSMeter(1.0)
    fps.FpsCap = 60

    glRenderer := &renderer.OpenglRenderer{
        WindowTitle : "GoEngine",
        WindowWidth : 900,
        WindowHeight : 700,
    }

    assetLib,err := assets.LoadAssetLibrary("TestAssets/demo.asset")
    if err != nil {
        panic(err)
    }

    //setup scenegraph

    geom := assetLib.GetGeometry("skybox")
    skyboxMat := assetLib.GetMaterial("skyboxMat")
    geom.Material = &skyboxMat
    geom.Material.LightingMode = renderer.MODE_UNLIT
    geom.CullBackface = false
    skyNode := renderer.CreateNode()
    skyNode.BucketType = renderer.BUCKET_OPAQUE
    skyNode.Add(&geom)
    skyNode.SetRotation( 1.57, vectorMath.Vector3{0,1,0} )
    skyNode.SetScale( vectorMath.Vector3{5000, 5000, 5000} )

    geomsphere := assetLib.GetGeometry("sphere")
    sphereMat := assetLib.GetMaterial("sphereMat")
    geomsphere.Material = &sphereMat
    boxNode2 := renderer.CreateNode()
    boxNode2.Add(&geomsphere)

    fireMat := assets.CreateMaterial(assetLib.GetImage("fire"), nil, nil, nil)
    fireMat.LightingMode = renderer.MODE_UNLIT
    firesprite := effects.CreateSprite( 36, 6, 6, &fireMat )
    firespriteNode := renderer.CreateNode()
    firespriteNode.Add(&firesprite)

    smokeMat := assets.CreateMaterial(assetLib.GetImage("smoke"), nil, nil, nil)
    smokeMat.LightingMode = renderer.MODE_UNLIT
    smokesprite := effects.CreateSprite( 64, 8, 8, &smokeMat )
    //smoke particle effect
    smokeParticles := effects.CreateParticleSystem( effects.ParticleSettings{
    	MaxParticles: 100,
		ParticleEmitRate: 10,
		Sprite: smokesprite,
		FaceCamera: true,
		MaxLife: 5.0,
		MinLife: 7.0,
		StartSize: vectorMath.Vector3{0.4, 0.4, 0.4},
		EndSize: vectorMath.Vector3{2.4, 2.4, 2.4},
		StartColor: color.NRGBA{254, 254, 254, 254},
		EndColor: color.NRGBA{254, 254, 254, 0},
		MinTranslation: vectorMath.Vector3{-0.2, -0.2, -0.2},
		MaxTranslation: vectorMath.Vector3{0.2, 0.2, 0.2},
		MaxStartVelocity: vectorMath.Vector3{-0.2, 0.3, 0.2},
		MinStartVelocity: vectorMath.Vector3{-0.2, 0.5, 0.2},
		Acceleration: vectorMath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity: vectorMath.IdentityQuaternion(),
		MinAngularVelocity: vectorMath.IdentityQuaternion(),
		MaxRotationVelocity: 0.0,
		MinRotationVelocity: 0.0,
    })

    explosionMat := assets.CreateMaterial(assetLib.GetImage("explosion"), nil, nil, nil)
    explosionMat.LightingMode = renderer.MODE_UNLIT
    explosionsprite := effects.CreateSprite( 36, 6, 6, &explosionMat )
    explosionspriteNode := renderer.CreateNode()
    explosionspriteNode.Add(&explosionsprite)
    explosionspriteNode.SetTranslation( vectorMath.Vector3{2,0,0} )

    sceneGraph := renderer.CreateSceneGraph()
    sceneGraph.Add(&skyNode)
    sceneGraph.Add(&boxNode2)
    sceneGraph.Add(&firespriteNode)
    sceneGraph.Add(&smokeParticles.Node)
    sceneGraph.Add(&explosionspriteNode)

    i := -45.0

    glRenderer.Init = func(){
        //setup reflection map
        cubeMap := renderer.CreateCubemap( assetLib.GetMaterial("skyboxMat").Diffuse );
        glRenderer.ReflectionMap( *cubeMap )
    }

    glRenderer.Update = func(){
        fps.UpdateFPSMeter()
        i = i + 0.11
        if i > 180 {
            i = -45
        }
        sine := math.Sin((float64)(i/26))
        cosine := math.Cos((float64)(i/26))

        boxNode2.SetTranslation( vectorMath.Vector3{1, 2, i} )
        //look at the box
        cameraLocation := vectorMath.Vector3{5*cosine,3*sine,5*sine}
        glRenderer.Camera( cameraLocation, vectorMath.Vector3{0,0,0}, vectorMath.Vector3{0,1,0} )

        glRenderer.CreateLight( 5,5,5, 100,100,100, 100,100,100, false, vectorMath.Vector3{1, 2, (float64)(i)}, 1 )

        //face the camera
        firespriteNode.SetFacing( 3.14, glRenderer.CameraLocation().Subtract(firespriteNode.Translation).Normalize(), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,-1} )
        explosionspriteNode.SetFacing( 3.14, glRenderer.CameraLocation().Subtract(explosionspriteNode.Translation).Normalize(), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,-1} )

        firesprite.NextFrame()
        explosionsprite.NextFrame()
        smokeParticles.Update(0.018, glRenderer)
    }

    glRenderer.Render = func(){
        sceneGraph.RenderScene(glRenderer)
    }

    glRenderer.Start()
}